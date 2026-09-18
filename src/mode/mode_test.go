package mode

import (
	"errors"
	"os"
	"strings"
	"testing"
)

// resetMode restores the package-level mode to Production after each test
// so tests don't leak state into one another (currentMode is global).
func resetMode(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		_ = Set("production")
	})
}

// TestParseMode covers every accepted alias plus the invalid-input error path.
func TestParseMode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Mode
		wantErr bool
	}{
		{name: "dev alias", input: "dev", want: Development},
		{name: "development", input: "development", want: Development},
		{name: "prod alias", input: "prod", want: Production},
		{name: "production", input: "production", want: Production},
		{name: "empty string", input: "", wantErr: true},
		{name: "garbage", input: "not-a-mode", wantErr: true},
		{name: "case sensitive", input: "PRODUCTION", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMode(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseMode(%q) expected error, got nil", tt.input)
				}
				if !strings.Contains(err.Error(), tt.input) {
					t.Errorf("error %q does not mention input %q", err.Error(), tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseMode(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ParseMode(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestSetAndGet covers the happy path (valid mode persists) and the error
// path (invalid mode is rejected and does not change current state).
func TestSetAndGet(t *testing.T) {
	resetMode(t)

	if err := Set("development"); err != nil {
		t.Fatalf("Set(development) unexpected error: %v", err)
	}
	if Get() != Development {
		t.Errorf("Get() = %v, want %v", Get(), Development)
	}
	if !IsDevelopment() {
		t.Error("IsDevelopment() = false, want true")
	}
	if IsProduction() {
		t.Error("IsProduction() = true, want false")
	}

	// Invalid Set must not mutate the current mode.
	if err := Set("bogus"); err == nil {
		t.Fatal("Set(bogus) expected error, got nil")
	}
	if Get() != Development {
		t.Errorf("Get() after failed Set = %v, want unchanged %v", Get(), Development)
	}

	if err := Set("production"); err != nil {
		t.Fatalf("Set(production) unexpected error: %v", err)
	}
	if !IsProduction() {
		t.Error("IsProduction() = false, want true")
	}
}

// TestInitialize covers CLI-flag priority, MODE env var priority, and the
// default-to-production fallback.
func TestInitialize(t *testing.T) {
	resetMode(t)

	t.Run("cli flag wins", func(t *testing.T) {
		os.Setenv("MODE", "production")
		defer os.Unsetenv("MODE")

		if err := Initialize("development"); err != nil {
			t.Fatalf("Initialize(development) unexpected error: %v", err)
		}
		if Get() != Development {
			t.Errorf("Get() = %v, want %v (CLI flag should win)", Get(), Development)
		}
	})

	t.Run("env var used when no cli flag", func(t *testing.T) {
		os.Setenv("MODE", "development")
		defer os.Unsetenv("MODE")

		if err := Initialize(""); err != nil {
			t.Fatalf("Initialize(\"\") unexpected error: %v", err)
		}
		if Get() != Development {
			t.Errorf("Get() = %v, want %v (MODE env var should be used)", Get(), Development)
		}
	})

	t.Run("defaults to production", func(t *testing.T) {
		os.Unsetenv("MODE")

		if err := Initialize(""); err != nil {
			t.Fatalf("Initialize(\"\") unexpected error: %v", err)
		}
		if Get() != Production {
			t.Errorf("Get() = %v, want %v (default)", Get(), Production)
		}
	})

	t.Run("invalid cli flag propagates error", func(t *testing.T) {
		if err := Initialize("nonsense"); err == nil {
			t.Fatal("Initialize(nonsense) expected error, got nil")
		}
	})
}

// TestGetErrorDetail covers nil error, development mode (stack trace
// included), and production mode (generic message only).
func TestGetErrorDetail(t *testing.T) {
	resetMode(t)

	if got := GetErrorDetail(nil); got != "" {
		t.Errorf("GetErrorDetail(nil) = %q, want empty string", got)
	}

	sampleErr := errors.New("boom")

	_ = Set("production")
	prodDetail := GetErrorDetail(sampleErr)
	if prodDetail != "An internal error occurred. Please try again later." {
		t.Errorf("GetErrorDetail() in production = %q, want generic message", prodDetail)
	}

	_ = Set("development")
	devDetail := GetErrorDetail(sampleErr)
	if !strings.Contains(devDetail, "boom") {
		t.Errorf("GetErrorDetail() in development = %q, want to contain error message", devDetail)
	}
	if !strings.Contains(devDetail, "Stack trace:") {
		t.Errorf("GetErrorDetail() in development = %q, want stack trace", devDetail)
	}
}

// TestModeDependentHelpers covers every helper that branches on the current
// mode, in both states.
func TestModeDependentHelpers(t *testing.T) {
	resetMode(t)

	_ = Set("development")
	if !ShouldShowDebugEndpoints() {
		t.Error("ShouldShowDebugEndpoints() in dev = false, want true")
	}
	if ShouldCacheTemplates() {
		t.Error("ShouldCacheTemplates() in dev = true, want false")
	}
	if ShouldCacheStaticFiles() {
		t.Error("ShouldCacheStaticFiles() in dev = true, want false")
	}
	if !ShouldEnableAutoReload() {
		t.Error("ShouldEnableAutoReload() in dev = false, want true")
	}
	if !ShouldEnableProfiling() {
		t.Error("ShouldEnableProfiling() in dev = false, want true")
	}
	if GetLogLevel() != "debug" {
		t.Errorf("GetLogLevel() in dev = %q, want debug", GetLogLevel())
	}
	devHeaders := GetCacheHeaders()
	if devHeaders["Cache-Control"] != "no-cache, no-store, must-revalidate" {
		t.Errorf("GetCacheHeaders() in dev = %v, want no-cache header", devHeaders)
	}
	if !strings.Contains(GetPanicRecoveryDetail("panic!"), "Stack trace:") {
		t.Error("GetPanicRecoveryDetail() in dev should include stack trace")
	}

	_ = Set("production")
	if ShouldShowDebugEndpoints() {
		t.Error("ShouldShowDebugEndpoints() in prod = true, want false")
	}
	if !ShouldCacheTemplates() {
		t.Error("ShouldCacheTemplates() in prod = false, want true")
	}
	if !ShouldCacheStaticFiles() {
		t.Error("ShouldCacheStaticFiles() in prod = false, want true")
	}
	if ShouldEnableAutoReload() {
		t.Error("ShouldEnableAutoReload() in prod = true, want false")
	}
	if ShouldEnableProfiling() {
		t.Error("ShouldEnableProfiling() in prod = true, want false")
	}
	if GetLogLevel() != "info" {
		t.Errorf("GetLogLevel() in prod = %q, want info", GetLogLevel())
	}
	prodHeaders := GetCacheHeaders()
	if prodHeaders["Cache-Control"] != "public, max-age=31536000, immutable" {
		t.Errorf("GetCacheHeaders() in prod = %v, want long-lived cache header", prodHeaders)
	}
	if GetPanicRecoveryDetail("panic!") != "Internal Server Error" {
		t.Errorf("GetPanicRecoveryDetail() in prod = %q, want generic message", GetPanicRecoveryDetail("panic!"))
	}
}

// TestModeString covers the Stringer implementation.
func TestModeString(t *testing.T) {
	if Development.String() != "development" {
		t.Errorf("Development.String() = %q, want development", Development.String())
	}
	if Production.String() != "production" {
		t.Errorf("Production.String() = %q, want production", Production.String())
	}
}
