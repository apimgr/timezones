package config

import (
	"os"
	"path/filepath"
	"testing"
)

// resetGlobalState clears the package-level current/configPath globals after
// each test so tests don't leak state into one another.
func resetGlobalState(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		mu.Lock()
		current = nil
		configPath = ""
		mu.Unlock()
	})
}

// TestDefaultConfig spot-checks representative fields across every nested
// struct to catch accidental default-value regressions.
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Server.Address != "0.0.0.0" {
		t.Errorf("Server.Address = %q, want %q", cfg.Server.Address, "0.0.0.0")
	}
	if cfg.Server.Mode != "production" {
		t.Errorf("Server.Mode = %q, want %q", cfg.Server.Mode, "production")
	}
	if cfg.Server.UpdateBranch != "stable" {
		t.Errorf("Server.UpdateBranch = %q, want %q", cfg.Server.UpdateBranch, "stable")
	}
	if cfg.Server.Metrics.Enabled {
		t.Error("Server.Metrics.Enabled = true, want false")
	}
	if cfg.Server.Metrics.Endpoint != "/metrics" {
		t.Errorf("Server.Metrics.Endpoint = %q, want %q", cfg.Server.Metrics.Endpoint, "/metrics")
	}
	if cfg.Server.Logging.Level != "info" {
		t.Errorf("Server.Logging.Level = %q, want %q", cfg.Server.Logging.Level, "info")
	}
	if cfg.WebUI.Theme != "dark" {
		t.Errorf("WebUI.Theme = %q, want %q", cfg.WebUI.Theme, "dark")
	}
	if len(cfg.WebRobots.Allow) != 2 {
		t.Errorf("WebRobots.Allow = %v, want 2 entries", cfg.WebRobots.Allow)
	}
	if cfg.WebSecurity.CORS != "*" {
		t.Errorf("WebSecurity.CORS = %q, want %q", cfg.WebSecurity.CORS, "*")
	}
}

// TestLoad covers: file doesn't exist (creates default + returns it), file
// exists with partial YAML (merges over defaults), and malformed YAML
// (returns an error).
func TestLoad(t *testing.T) {
	t.Run("creates default when missing", func(t *testing.T) {
		resetGlobalState(t)
		dir := t.TempDir()
		path := filepath.Join(dir, "server.yml")

		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("Load() unexpected error: %v", err)
		}
		if cfg.Server.Address != "0.0.0.0" {
			t.Errorf("Load() returned Server.Address = %q, want default", cfg.Server.Address)
		}
		if _, err := os.Stat(path); err != nil {
			t.Errorf("Load() did not create config file: %v", err)
		}
	})

	t.Run("loads existing partial YAML over defaults", func(t *testing.T) {
		resetGlobalState(t)
		dir := t.TempDir()
		path := filepath.Join(dir, "server.yml")
		partial := "server:\n  port: \"9999\"\n  address: \"127.0.0.1\"\n"
		if err := os.WriteFile(path, []byte(partial), 0644); err != nil {
			t.Fatalf("failed to write fixture config: %v", err)
		}

		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("Load() unexpected error: %v", err)
		}
		if cfg.Server.Port != "9999" {
			t.Errorf("Server.Port = %q, want %q", cfg.Server.Port, "9999")
		}
		if cfg.Server.Address != "127.0.0.1" {
			t.Errorf("Server.Address = %q, want %q", cfg.Server.Address, "127.0.0.1")
		}
		// Fields absent from the partial YAML must still carry defaults.
		if cfg.Server.Mode != "production" {
			t.Errorf("Server.Mode = %q, want default %q", cfg.Server.Mode, "production")
		}
	})

	t.Run("malformed YAML returns error", func(t *testing.T) {
		resetGlobalState(t)
		dir := t.TempDir()
		path := filepath.Join(dir, "server.yml")
		if err := os.WriteFile(path, []byte("server: [this is not valid: yaml"), 0644); err != nil {
			t.Fatalf("failed to write fixture config: %v", err)
		}

		if _, err := Load(path); err == nil {
			t.Fatal("Load() expected error for malformed YAML, got nil")
		}
	})
}

// TestGet covers the nil-current fallback and the loaded-current passthrough.
func TestGet(t *testing.T) {
	t.Run("nil current returns default", func(t *testing.T) {
		resetGlobalState(t)
		cfg := Get()
		if cfg.Server.Address != "0.0.0.0" {
			t.Errorf("Get() with no current = %+v, want default", cfg)
		}
	})

	t.Run("returns loaded current", func(t *testing.T) {
		resetGlobalState(t)
		dir := t.TempDir()
		path := filepath.Join(dir, "server.yml")
		if _, err := Load(path); err != nil {
			t.Fatalf("Load() unexpected error: %v", err)
		}

		cfg := Get()
		if cfg == nil {
			t.Fatal("Get() returned nil after Load()")
		}
	})
}

// TestSave covers the no-config-loaded error path and the happy path
// (writes a file, and the values it does serialize round-trip). It also
// documents a known bug: generateConfigYAML does not serialize
// Server.Mode or Server.UpdateBranch, so those fields silently revert to
// defaults after a Save+Load round-trip. This is existing behavior, not
// fixed here per instruction.
func TestSave(t *testing.T) {
	t.Run("no config loaded returns error", func(t *testing.T) {
		resetGlobalState(t)
		if err := Save(); err == nil {
			t.Fatal("Save() expected error when no config loaded, got nil")
		}
	})

	t.Run("writes file and round-trips serialized fields", func(t *testing.T) {
		resetGlobalState(t)
		dir := t.TempDir()
		path := filepath.Join(dir, "server.yml")

		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("Load() unexpected error: %v", err)
		}
		cfg.Server.Port = "1234"
		cfg.Server.FQDN = "example.com"
		cfg.WebUI.Theme = "light"

		if err := Save(); err != nil {
			t.Fatalf("Save() unexpected error: %v", err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read saved config: %v", err)
		}
		content := string(data)
		if !contains(content, `port: "1234"`) {
			t.Errorf("saved config missing port: %s", content)
		}
		if !contains(content, `fqdn: "example.com"`) {
			t.Errorf("saved config missing fqdn: %s", content)
		}
		if !contains(content, `theme: "light"`) {
			t.Errorf("saved config missing theme: %s", content)
		}

		// Known bug: Mode/UpdateBranch are silently dropped by
		// generateConfigYAML, so they never appear in the file.
		if contains(content, "update_branch") {
			t.Error("unexpected: update_branch is now serialized (bug may have been fixed; update this test)")
		}
	})
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}

// TestGetTheme verifies the passthrough to WebUI.Theme.
func TestGetTheme(t *testing.T) {
	resetGlobalState(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "server.yml")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	cfg.WebUI.Theme = "custom"

	if got := GetTheme(); got != "custom" {
		t.Errorf("GetTheme() = %q, want %q", got, "custom")
	}
}

// TestGetCORS covers the empty-CORS fallback to "*" and the explicit
// passthrough case.
func TestGetCORS(t *testing.T) {
	t.Run("empty falls back to wildcard", func(t *testing.T) {
		resetGlobalState(t)
		dir := t.TempDir()
		path := filepath.Join(dir, "server.yml")
		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("Load() unexpected error: %v", err)
		}
		cfg.WebSecurity.CORS = ""

		if got := GetCORS(); got != "*" {
			t.Errorf("GetCORS() = %q, want %q", got, "*")
		}
	})

	t.Run("explicit value passes through", func(t *testing.T) {
		resetGlobalState(t)
		dir := t.TempDir()
		path := filepath.Join(dir, "server.yml")
		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("Load() unexpected error: %v", err)
		}
		cfg.WebSecurity.CORS = "https://example.com"

		if got := GetCORS(); got != "https://example.com" {
			t.Errorf("GetCORS() = %q, want %q", got, "https://example.com")
		}
	})
}

// TestFormatStringSlice covers the empty-slice and populated-slice branches.
func TestFormatStringSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{name: "empty", input: []string{}, want: "[]"},
		{name: "nil", input: nil, want: "[]"},
		{name: "single", input: []string{"a"}, want: `["a"]`},
		{name: "multiple", input: []string{"a", "b"}, want: `["a", "b"]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatStringSlice(tt.input); got != tt.want {
				t.Errorf("formatStringSlice(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestMigrateYamlToYml covers: no .yml suffix (no-op), .yml already exists
// (no-op), .yaml sibling exists and gets renamed, and neither exists (no-op).
func TestMigrateYamlToYml(t *testing.T) {
	t.Run("non-.yml path is a no-op", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "server.yaml")
		if err := migrateYamlToYml(path); err != nil {
			t.Fatalf("migrateYamlToYml() unexpected error: %v", err)
		}
	})

	t.Run(".yml already exists is a no-op", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "server.yml")
		if err := os.WriteFile(path, []byte("existing"), 0644); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}
		yamlPath := filepath.Join(dir, "server.yaml")
		if err := os.WriteFile(yamlPath, []byte("old"), 0644); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}

		if err := migrateYamlToYml(path); err != nil {
			t.Fatalf("migrateYamlToYml() unexpected error: %v", err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read %s: %v", path, err)
		}
		if string(data) != "existing" {
			t.Errorf(".yml content = %q, want unchanged %q", data, "existing")
		}
	})

	t.Run(".yaml sibling is renamed to .yml", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "server.yml")
		yamlPath := filepath.Join(dir, "server.yaml")
		if err := os.WriteFile(yamlPath, []byte("migrated content"), 0644); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}

		if err := migrateYamlToYml(path); err != nil {
			t.Fatalf("migrateYamlToYml() unexpected error: %v", err)
		}
		if _, err := os.Stat(yamlPath); !os.IsNotExist(err) {
			t.Error(".yaml file should no longer exist after migration")
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read migrated .yml: %v", err)
		}
		if string(data) != "migrated content" {
			t.Errorf(".yml content = %q, want %q", data, "migrated content")
		}
	})

	t.Run("neither exists is a no-op", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "server.yml")
		if err := migrateYamlToYml(path); err != nil {
			t.Fatalf("migrateYamlToYml() unexpected error: %v", err)
		}
	})
}
