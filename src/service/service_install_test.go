package service

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// withTestSysRoot redirects every absolute path this package writes to
// (via withRoot) into a fresh t.TempDir(), and restores the previous value
// on cleanup. This lets install/uninstall logic run for real without
// touching the actual filesystem.
func withTestSysRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	orig := sysRoot
	sysRoot = root
	t.Cleanup(func() { sysRoot = orig })
	return root
}

// stubRunCommand replaces runCommand with a fake that records every
// invocation instead of executing a real system command, and restores the
// original on cleanup. When err is non-nil, every call returns it.
func stubRunCommand(t *testing.T, err error) *[]string {
	t.Helper()
	calls := []string{}
	orig := runCommand
	runCommand = func(name string, args ...string) error {
		calls = append(calls, strings.TrimSpace(name+" "+strings.Join(args, " ")))
		return err
	}
	t.Cleanup(func() { runCommand = orig })
	return &calls
}

// stubDetect forces detectFn to report st, and restores the original on
// cleanup, so each switch branch in Install/Uninstall/Start/Stop/Restart/
// Reload can be exercised deterministically.
func stubDetect(t *testing.T, st ServiceType) {
	t.Helper()
	orig := detectFn
	detectFn = func() ServiceType { return st }
	t.Cleanup(func() { detectFn = orig })
}

// TestWithRoot covers both branches directly: unchanged when sysRoot is
// empty (the default, production behavior) and joined when it is set.
func TestWithRoot(t *testing.T) {
	orig := sysRoot
	defer func() { sysRoot = orig }()

	sysRoot = ""
	if got := withRoot("/etc/foo"); got != "/etc/foo" {
		t.Errorf("withRoot with empty sysRoot = %q, want %q", got, "/etc/foo")
	}

	sysRoot = "/tmp/root"
	want := filepath.Join("/tmp/root", "/etc/foo")
	if got := withRoot("/etc/foo"); got != want {
		t.Errorf("withRoot with sysRoot set = %q, want %q", got, want)
	}
}

// TestGetBinaryPathWithSysRoot verifies GetBinaryPath honors sysRoot.
func TestGetBinaryPathWithSysRoot(t *testing.T) {
	root := withTestSysRoot(t)
	want := filepath.Join(root, "/usr/local/bin/timezones")
	if got := GetBinaryPath(); got != want {
		t.Errorf("GetBinaryPath() = %q, want %q", got, want)
	}
}

// TestInstallUninstallUnsupported covers the default (unsupported service
// manager) branch of Install and Uninstall without touching the
// filesystem or spawning any command.
func TestInstallUninstallUnsupported(t *testing.T) {
	stubDetect(t, ServiceUnknown)

	if err := Install(); err == nil {
		t.Fatal("Install() with ServiceUnknown: expected error, got nil")
	}
	if err := Uninstall(); err == nil {
		t.Fatal("Uninstall() with ServiceUnknown: expected error, got nil")
	}
}

// TestInstallSystemd verifies the generated unit file content, the created
// data/log/config directories, and that daemon-reload/enable are invoked
// via runCommand rather than a real systemctl.
func TestInstallSystemd(t *testing.T) {
	root := withTestSysRoot(t)
	calls := stubRunCommand(t, nil)
	stubDetect(t, ServiceSystemd)

	// The real systemd unit directory always pre-exists on a host that has
	// systemd; recreate that precondition under the fake root.
	if err := os.MkdirAll(filepath.Join(root, "/etc/systemd/system"), 0755); err != nil {
		t.Fatalf("setup MkdirAll: %v", err)
	}

	if err := Install(); err != nil {
		t.Fatalf("Install() unexpected error: %v", err)
	}

	servicePath := filepath.Join(root, "/etc/systemd/system/timezones.service")
	content, err := os.ReadFile(servicePath)
	if err != nil {
		t.Fatalf("service file not written: %v", err)
	}
	binaryPath := filepath.Join(root, "/usr/local/bin/timezones")
	for _, want := range []string{
		"ExecStart=" + binaryPath,
		"ReadWritePaths=/var/lib/apimgr/timezones /var/log/apimgr/timezones /etc/apimgr/timezones",
		"[Install]",
	} {
		if !strings.Contains(string(content), want) {
			t.Errorf("service file missing %q\ngot:\n%s", want, content)
		}
	}

	for _, dir := range []string{
		"/var/lib/apimgr/timezones",
		"/var/log/apimgr/timezones",
		"/etc/apimgr/timezones",
	} {
		info, err := os.Stat(filepath.Join(root, dir))
		if err != nil || !info.IsDir() {
			t.Errorf("expected directory %s to exist: %v", dir, err)
		}
	}

	wantCalls := []string{"systemctl daemon-reload", "systemctl enable timezones"}
	if strings.Join(*calls, "|") != strings.Join(wantCalls, "|") {
		t.Errorf("runCommand calls = %v, want %v", *calls, wantCalls)
	}

	// Uninstall should remove the unit file and issue stop/disable/reload.
	*calls = nil
	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() unexpected error: %v", err)
	}
	if _, err := os.Stat(servicePath); !os.IsNotExist(err) {
		t.Errorf("expected service file removed, stat err = %v", err)
	}
	wantUninstall := []string{"systemctl stop timezones", "systemctl disable timezones", "systemctl daemon-reload"}
	if strings.Join(*calls, "|") != strings.Join(wantUninstall, "|") {
		t.Errorf("uninstall runCommand calls = %v, want %v", *calls, wantUninstall)
	}

	// Uninstalling again (file already gone) must not error.
	if err := Uninstall(); err != nil {
		t.Errorf("second Uninstall() unexpected error: %v", err)
	}
}

// TestInstallSystemdReloadError verifies the failure path is surfaced when
// the systemd reload command fails.
func TestInstallSystemdReloadError(t *testing.T) {
	root := withTestSysRoot(t)
	stubRunCommand(t, errors.New("boom"))
	stubDetect(t, ServiceSystemd)

	if err := os.MkdirAll(filepath.Join(root, "/etc/systemd/system"), 0755); err != nil {
		t.Fatalf("setup MkdirAll: %v", err)
	}

	err := Install()
	if err == nil {
		t.Fatal("Install() expected error when systemctl daemon-reload fails")
	}
	if !strings.Contains(err.Error(), "failed to reload systemd") {
		t.Errorf("Install() error = %v, want it to mention reload failure", err)
	}
}

// TestInstallRunit verifies the run script, log run script, and service
// symlink content/placement.
func TestInstallRunit(t *testing.T) {
	root := withTestSysRoot(t)
	calls := stubRunCommand(t, nil)
	stubDetect(t, ServiceRunit)

	// os.Symlink requires its parent directory to already exist, and the
	// real /var/service directory always pre-exists on a runit host;
	// recreate that precondition under the fake root. (installRunit
	// intentionally ignores the symlink error, matching upstream.)
	if err := os.MkdirAll(filepath.Join(root, "/var/service"), 0755); err != nil {
		t.Fatalf("setup MkdirAll: %v", err)
	}

	if err := Install(); err != nil {
		t.Fatalf("Install() unexpected error: %v", err)
	}

	svDir := filepath.Join(root, "/etc/sv/timezones")
	binaryPath := filepath.Join(root, "/usr/local/bin/timezones")

	runScript, err := os.ReadFile(filepath.Join(svDir, "run"))
	if err != nil {
		t.Fatalf("run script not written: %v", err)
	}
	if !strings.Contains(string(runScript), "exec "+binaryPath) {
		t.Errorf("run script = %q, want it to exec %q", runScript, binaryPath)
	}

	logScript, err := os.ReadFile(filepath.Join(svDir, "log", "run"))
	if err != nil {
		t.Fatalf("log run script not written: %v", err)
	}
	if !strings.Contains(string(logScript), "svlogd") {
		t.Errorf("log run script = %q, want it to reference svlogd", logScript)
	}

	linkPath := filepath.Join(root, "/var/service/timezones")
	target, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("expected symlink at %s: %v", linkPath, err)
	}
	if target != svDir {
		t.Errorf("symlink target = %q, want %q", target, svDir)
	}

	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() unexpected error: %v", err)
	}
	if _, err := os.Stat(svDir); !os.IsNotExist(err) {
		t.Errorf("expected sv dir removed, stat err = %v", err)
	}
	if _, err := os.Lstat(linkPath); !os.IsNotExist(err) {
		t.Errorf("expected symlink removed, stat err = %v", err)
	}
	if len(*calls) == 0 || (*calls)[len(*calls)-1] != "sv stop timezones" {
		t.Errorf("expected a final 'sv stop timezones' call, got %v", *calls)
	}
}

// TestInstallLaunchd verifies the generated plist content and directories.
func TestInstallLaunchd(t *testing.T) {
	root := withTestSysRoot(t)
	stubRunCommand(t, nil)
	stubDetect(t, ServiceLaunchd)

	// The real /Library/LaunchDaemons directory always pre-exists on
	// macOS; recreate that precondition under the fake root.
	if err := os.MkdirAll(filepath.Join(root, "/Library/LaunchDaemons"), 0755); err != nil {
		t.Fatalf("setup MkdirAll: %v", err)
	}

	if err := Install(); err != nil {
		t.Fatalf("Install() unexpected error: %v", err)
	}

	plistPath := filepath.Join(root, "/Library/LaunchDaemons/com.apimgr.timezones.plist")
	content, err := os.ReadFile(plistPath)
	if err != nil {
		t.Fatalf("plist not written: %v", err)
	}
	binaryPath := filepath.Join(root, "/usr/local/bin/timezones")
	for _, want := range []string{
		"<string>com.apimgr.timezones</string>",
		"<string>" + binaryPath + "</string>",
		"error.log",
		"output.log",
	} {
		if !strings.Contains(string(content), want) {
			t.Errorf("plist missing %q\ngot:\n%s", want, content)
		}
	}

	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() unexpected error: %v", err)
	}
	if _, err := os.Stat(plistPath); !os.IsNotExist(err) {
		t.Errorf("expected plist removed, stat err = %v", err)
	}

	// Uninstalling again (already removed) must not error.
	if err := Uninstall(); err != nil {
		t.Errorf("second Uninstall() unexpected error: %v", err)
	}
}

// TestInstallBSDRC verifies the generated rc.d script content.
func TestInstallBSDRC(t *testing.T) {
	root := withTestSysRoot(t)
	calls := stubRunCommand(t, nil)
	stubDetect(t, ServiceBSDRC)

	// The real /usr/local/etc/rc.d directory always pre-exists on a BSD
	// host; recreate that precondition under the fake root.
	if err := os.MkdirAll(filepath.Join(root, "/usr/local/etc/rc.d"), 0755); err != nil {
		t.Fatalf("setup MkdirAll: %v", err)
	}

	if err := Install(); err != nil {
		t.Fatalf("Install() unexpected error: %v", err)
	}

	rcPath := filepath.Join(root, "/usr/local/etc/rc.d/timezones")
	content, err := os.ReadFile(rcPath)
	if err != nil {
		t.Fatalf("rc.d script not written: %v", err)
	}
	binaryPath := filepath.Join(root, "/usr/local/bin/timezones")
	for _, want := range []string{
		`name="timezones"`,
		`command="` + binaryPath + `"`,
		"PROVIDE: timezones",
	} {
		if !strings.Contains(string(content), want) {
			t.Errorf("rc.d script missing %q\ngot:\n%s", want, content)
		}
	}

	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() unexpected error: %v", err)
	}
	if _, err := os.Stat(rcPath); !os.IsNotExist(err) {
		t.Errorf("expected rc.d script removed, stat err = %v", err)
	}
	if len(*calls) == 0 || (*calls)[len(*calls)-1] != "service timezones stop" {
		t.Errorf("expected a final 'service timezones stop' call, got %v", *calls)
	}
}

// TestInstallWindows verifies the Windows install/uninstall path invokes
// sc.exe (via the stubbed runCommand) with the expected arguments and
// copies the binary into place, and that a create failure is surfaced.
func TestInstallWindows(t *testing.T) {
	withTestSysRoot(t)
	calls := stubRunCommand(t, nil)
	stubDetect(t, ServiceWindows)

	if err := Install(); err != nil {
		t.Fatalf("Install() unexpected error: %v", err)
	}
	if len(*calls) != 1 || !strings.HasPrefix((*calls)[0], "sc.exe create timezones") {
		t.Errorf("runCommand calls = %v, want a single sc.exe create call", *calls)
	}

	binaryPath := GetBinaryPath()
	if _, err := os.Stat(binaryPath); err != nil {
		t.Errorf("expected binary copied to %s: %v", binaryPath, err)
	}

	*calls = nil
	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() unexpected error: %v", err)
	}
	wantCalls := []string{"sc.exe stop timezones", "sc.exe delete timezones"}
	if strings.Join(*calls, "|") != strings.Join(wantCalls, "|") {
		t.Errorf("uninstall runCommand calls = %v, want %v", *calls, wantCalls)
	}
}

// TestInstallWindowsCreateError verifies Install surfaces an error when
// sc.exe create fails.
func TestInstallWindowsCreateError(t *testing.T) {
	withTestSysRoot(t)
	stubRunCommand(t, errors.New("access denied"))
	stubDetect(t, ServiceWindows)

	err := Install()
	if err == nil {
		t.Fatal("Install() expected error when sc.exe create fails")
	}
	if !strings.Contains(err.Error(), "failed to create Windows service") {
		t.Errorf("Install() error = %v, want it to mention Windows service creation failure", err)
	}
}

// TestUninstallWindowsDeleteError verifies Uninstall surfaces an error when
// sc.exe delete fails, even though the preceding stop call's error is
// ignored.
func TestUninstallWindowsDeleteError(t *testing.T) {
	withTestSysRoot(t)
	stubRunCommand(t, errors.New("not found"))
	stubDetect(t, ServiceWindows)

	err := Uninstall()
	if err == nil {
		t.Fatal("Uninstall() expected error when sc.exe delete fails")
	}
	if !strings.Contains(err.Error(), "failed to delete Windows service") {
		t.Errorf("Uninstall() error = %v, want it to mention delete failure", err)
	}
}

// TestStartStopRestartReloadDispatch exercises every branch of Start, Stop,
// Restart, and Reload for each known service manager, and the "unsupported"
// default branch, all via the stubbed runCommand/detectFn so no real
// service manager is ever invoked.
func TestStartStopRestartReloadDispatch(t *testing.T) {
	types := []ServiceType{ServiceSystemd, ServiceRunit, ServiceLaunchd, ServiceWindows, ServiceBSDRC, ServiceUnknown}

	for _, st := range types {
		t.Run(serviceTypeName(st), func(t *testing.T) {
			withTestSysRoot(t)
			stubRunCommand(t, nil)
			stubDetect(t, st)

			wantErr := st == ServiceUnknown

			if err := Start(); (err != nil) != wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, wantErr)
			}
			if err := Stop(); (err != nil) != wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, wantErr)
			}
			if err := Restart(); (err != nil) != wantErr {
				t.Errorf("Restart() error = %v, wantErr %v", err, wantErr)
			}
			if err := Reload(); (err != nil) != wantErr {
				t.Errorf("Reload() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}

// TestStartCommandErrorPropagates verifies a runCommand failure is
// returned by Start for a simple (non-launchd, non-windows) dispatch.
func TestStartCommandErrorPropagates(t *testing.T) {
	withTestSysRoot(t)
	stubRunCommand(t, errors.New("systemctl failed"))
	stubDetect(t, ServiceSystemd)

	if err := Start(); err == nil {
		t.Fatal("Start() expected error to propagate from runCommand")
	}
}

func serviceTypeName(st ServiceType) string {
	switch st {
	case ServiceSystemd:
		return "systemd"
	case ServiceRunit:
		return "runit"
	case ServiceLaunchd:
		return "launchd"
	case ServiceWindows:
		return "windows"
	case ServiceBSDRC:
		return "bsdrc"
	default:
		return "unknown"
	}
}
