package service

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestDetectServiceManager verifies a value from the known set is returned
// for the current runtime.GOOS, and never panics. The exact value is
// environment-dependent (host may or may not have systemd/runit), so we only
// assert it is one of the defined ServiceType constants.
func TestDetectServiceManager(t *testing.T) {
	got := DetectServiceManager()

	valid := map[ServiceType]bool{
		ServiceUnknown:  true,
		ServiceSystemd:  true,
		ServiceRunit:    true,
		ServiceLaunchd:  true,
		ServiceWindows:  true,
		ServiceBSDRC:    true,
	}
	if !valid[got] {
		t.Errorf("DetectServiceManager() = %v, want one of the defined ServiceType values", got)
	}

	switch runtime.GOOS {
	case "darwin":
		if got != ServiceLaunchd {
			t.Errorf("DetectServiceManager() on darwin = %v, want ServiceLaunchd", got)
		}
	case "windows":
		if got != ServiceWindows {
			t.Errorf("DetectServiceManager() on windows = %v, want ServiceWindows", got)
		}
	case "freebsd", "openbsd", "netbsd":
		if got != ServiceBSDRC {
			t.Errorf("DetectServiceManager() on %s = %v, want ServiceBSDRC", runtime.GOOS, got)
		}
	}
}

// TestGetBinaryPath verifies the expected path shape for the current OS.
func TestGetBinaryPath(t *testing.T) {
	got := GetBinaryPath()
	if got == "" {
		t.Fatal("GetBinaryPath() returned empty string")
	}

	if runtime.GOOS == "windows" {
		want := `C:\Program Files\apimgr\timezones\timezones.exe`
		if got != want {
			t.Errorf("GetBinaryPath() = %q, want %q", got, want)
		}
		return
	}

	want := "/usr/local/bin/timezones"
	if got != want {
		t.Errorf("GetBinaryPath() = %q, want %q", got, want)
	}
}

// TestCopyBinary covers copying a file into a new nested destination
// directory (verifying MkdirAll is applied), preserving content, and the
// error path when the source does not exist.
func TestCopyBinary(t *testing.T) {
	t.Run("copies content and creates parent dirs", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "src-binary")
		content := []byte("fake binary contents")
		if err := os.WriteFile(src, content, 0755); err != nil {
			t.Fatalf("WriteFile(src): %v", err)
		}

		dst := filepath.Join(dir, "nested", "deeper", "dst-binary")
		if err := copyBinary(src, dst); err != nil {
			t.Fatalf("copyBinary() unexpected error: %v", err)
		}

		got, err := os.ReadFile(dst)
		if err != nil {
			t.Fatalf("ReadFile(dst): %v", err)
		}
		if string(got) != string(content) {
			t.Errorf("copied content = %q, want %q", got, content)
		}
	})

	t.Run("error when source does not exist", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "does-not-exist")
		dst := filepath.Join(dir, "dst")

		if err := copyBinary(src, dst); err == nil {
			t.Fatal("copyBinary() expected error for missing source, got nil")
		}
	})
}

// TestTitleCase covers the empty string, a lowercase word, and a string that
// starts already capitalized.
func TestTitleCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty string", input: "", want: ""},
		{name: "lowercase word", input: "timezones", want: "Timezones"},
		{name: "already capitalized", input: "Timezones", want: "Timezones"},
		{name: "single character", input: "a", want: "A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := titleCase(tt.input); got != tt.want {
				t.Errorf("titleCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
