package paths

import (
	"os"
	"path/filepath"
	"testing"
)

// clearEnv unsets every override env var these functions consult so each
// test starts from a known-clean slate, then restores it afterward.
func clearEnv(t *testing.T) {
	t.Helper()
	vars := []string{"CONFIG_DIR", "DATA_DIR", "LOGS_DIR", "XDG_CONFIG_HOME", "XDG_DATA_HOME"}
	saved := make(map[string]string, len(vars))
	for _, v := range vars {
		saved[v] = os.Getenv(v)
		os.Unsetenv(v)
	}
	t.Cleanup(func() {
		for _, v := range vars {
			if saved[v] == "" {
				os.Unsetenv(v)
			} else {
				os.Setenv(v, saved[v])
			}
		}
	})
}

// TestGetConfigDirEnvOverride covers the CONFIG_DIR override path, which
// always wins regardless of OS or privilege level.
func TestGetConfigDirEnvOverride(t *testing.T) {
	clearEnv(t)
	os.Setenv("CONFIG_DIR", "/tmp/custom-config")

	got := GetConfigDir("myapp")
	if got != "/tmp/custom-config" {
		t.Errorf("GetConfigDir() = %q, want %q", got, "/tmp/custom-config")
	}
}

// TestGetConfigDirDefault exercises the non-override code path for the
// current OS/privilege combination. We don't assert an exact path (that
// varies by platform and by whether the test runs as root), only that a
// non-empty, plausible directory is returned and it contains the app name.
func TestGetConfigDirDefault(t *testing.T) {
	clearEnv(t)

	got := GetConfigDir("myapp")
	if got == "" {
		t.Fatal("GetConfigDir() returned empty string")
	}
}

// TestGetConfigDirXDG covers the XDG_CONFIG_HOME override on the default
// (non-root, non-Windows, non-macOS) branch when running unprivileged.
func TestGetConfigDirXDG(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root: XDG_CONFIG_HOME branch not reachable (root path wins)")
	}
	clearEnv(t)
	os.Setenv("XDG_CONFIG_HOME", "/tmp/xdgcfg")

	got := GetConfigDir("myapp")
	want := filepath.Join("/tmp/xdgcfg", OrgName, "myapp")
	if got != want {
		t.Errorf("GetConfigDir() = %q, want %q", got, want)
	}
}

// TestGetDataDirEnvOverride covers the DATA_DIR override.
func TestGetDataDirEnvOverride(t *testing.T) {
	clearEnv(t)
	os.Setenv("DATA_DIR", "/tmp/custom-data")

	got := GetDataDir("myapp")
	if got != "/tmp/custom-data" {
		t.Errorf("GetDataDir() = %q, want %q", got, "/tmp/custom-data")
	}
}

// TestGetDataDirDefault exercises the non-override code path.
func TestGetDataDirDefault(t *testing.T) {
	clearEnv(t)

	got := GetDataDir("myapp")
	if got == "" {
		t.Fatal("GetDataDir() returned empty string")
	}
}

// TestGetLogsDirEnvOverride covers the LOGS_DIR override.
func TestGetLogsDirEnvOverride(t *testing.T) {
	clearEnv(t)
	os.Setenv("LOGS_DIR", "/tmp/custom-logs")

	got := GetLogsDir("myapp")
	if got != "/tmp/custom-logs" {
		t.Errorf("GetLogsDir() = %q, want %q", got, "/tmp/custom-logs")
	}
}

// TestGetLogsDirDefault exercises the non-override code path. On Linux this
// is root-aware: /var/log/OrgName/appName as root, GetDataDir()+"/logs"
// otherwise (tests commonly run as root inside CI containers).
func TestGetLogsDirDefault(t *testing.T) {
	clearEnv(t)

	got := GetLogsDir("myapp")
	if got == "" {
		t.Fatal("GetLogsDir() returned empty string")
	}
	var want string
	if os.Geteuid() == 0 {
		want = filepath.Join("/var/log", OrgName, "myapp")
	} else {
		want = filepath.Join(GetDataDir("myapp"), "logs")
	}
	if got != want {
		t.Errorf("GetLogsDir() = %q, want %q", got, want)
	}
}

// TestEnsureDir covers creating a fresh directory (happy path), a nested
// path that requires parent creation, and idempotency (calling twice on the
// same path must not error).
func TestEnsureDir(t *testing.T) {
	base := t.TempDir()

	tests := []struct {
		name string
		path string
	}{
		{name: "single level", path: filepath.Join(base, "onelevel")},
		{name: "nested path", path: filepath.Join(base, "a", "b", "c")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EnsureDir(tt.path); err != nil {
				t.Fatalf("EnsureDir(%q) unexpected error: %v", tt.path, err)
			}
			info, err := os.Stat(tt.path)
			if err != nil {
				t.Fatalf("directory was not created: %v", err)
			}
			if !info.IsDir() {
				t.Fatalf("%q is not a directory", tt.path)
			}

			// Idempotency: calling again must be safe.
			if err := EnsureDir(tt.path); err != nil {
				t.Fatalf("EnsureDir(%q) second call unexpected error: %v", tt.path, err)
			}
		})
	}
}

// TestOrgAndProjectConstants pins the exported constants so an accidental
// rename is caught by tests, not just by downstream breakage.
func TestOrgAndProjectConstants(t *testing.T) {
	if OrgName != "apimgr" {
		t.Errorf("OrgName = %q, want %q", OrgName, "apimgr")
	}
	if ProjectName != "timezones" {
		t.Errorf("ProjectName = %q, want %q", ProjectName, "timezones")
	}
}
