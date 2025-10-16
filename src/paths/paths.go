package paths

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetConfigDir returns the OS-specific configuration directory
func GetConfigDir(appName string) string {
	// Check environment variable first
	if configDir := os.Getenv("CONFIG_DIR"); configDir != "" {
		return configDir
	}

	var baseDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%\AppName
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			baseDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
		return filepath.Join(baseDir, capitalize(appName))

	case "darwin":
		// macOS: ~/Library/Application Support/AppName
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, "Library", "Application Support", capitalize(appName))

	default:
		// Linux/Unix: Check if running as root
		if os.Geteuid() == 0 {
			// Root user: /etc/appname
			return filepath.Join("/etc", appName)
		}
		// Regular user: ~/.config/appname
		if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
			return filepath.Join(xdgConfig, appName)
		}
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".config", appName)
	}
}

// GetDataDir returns the OS-specific data directory
func GetDataDir(appName string) string {
	// Check environment variable first
	if dataDir := os.Getenv("DATA_DIR"); dataDir != "" {
		return dataDir
	}

	var baseDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: %LOCALAPPDATA%\AppName
		baseDir = os.Getenv("LOCALAPPDATA")
		if baseDir == "" {
			baseDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
		return filepath.Join(baseDir, capitalize(appName))

	case "darwin":
		// macOS: ~/Library/Application Support/AppName
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, "Library", "Application Support", capitalize(appName))

	default:
		// Linux/Unix: Check if running as root
		if os.Geteuid() == 0 {
			// Root user: /var/lib/appname
			return filepath.Join("/var/lib", appName)
		}
		// Regular user: ~/.local/share/appname
		if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
			return filepath.Join(xdgData, appName)
		}
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".local", "share", appName)
	}
}

// GetLogsDir returns the OS-specific logs directory
func GetLogsDir(appName string) string {
	// Check environment variable first
	if logsDir := os.Getenv("LOGS_DIR"); logsDir != "" {
		return logsDir
	}

	switch runtime.GOOS {
	case "windows":
		// Windows: %LOCALAPPDATA%\AppName\logs
		return filepath.Join(GetDataDir(appName), "logs")

	case "darwin":
		// macOS: ~/Library/Logs/AppName
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, "Library", "Logs", capitalize(appName))

	default:
		// Linux/Unix: Check if running as root
		if os.Geteuid() == 0 {
			// Root user: /var/log/appname
			return filepath.Join("/var/log", appName)
		}
		// Regular user: ~/.local/share/appname/logs
		return filepath.Join(GetDataDir(appName), "logs")
	}
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// capitalize returns the string with first letter capitalized
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
