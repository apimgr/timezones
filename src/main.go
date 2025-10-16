package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/casjay/timezones/src/database"
	"github.com/casjay/timezones/src/paths"
	"github.com/casjay/timezones/src/server"
	"github.com/casjay/timezones/src/timezones"
)

// Embed timezones JSON data
//go:embed data/timezones.json
var timezonesJSON []byte

// Build information (set via ldflags)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Command-line flags
	var (
		port         = flag.String("port", "", "Port to listen on (default: 8080)")
		address      = flag.String("address", "", "Address to bind to (default: 0.0.0.0)")
		configDir    = flag.String("config", "", "Configuration directory")
		dataDir      = flag.String("data", "", "Data directory")
		logsDir      = flag.String("logs", "", "Logs directory")
		adminUser    = flag.String("admin-user", "", "Admin username (for first-time setup)")
		adminPass    = flag.String("admin-password", "", "Admin password (for first-time setup)")
		showVersion  = flag.Bool("version", false, "Show version information")
		showStatus   = flag.Bool("status", false, "Show server status (for health checks)")
	)
	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("%s\n", Version)
		return
	}

	// Status check (for health checks)
	if *showStatus {
		os.Exit(0)
	}

	log.Printf("Starting Timezones API v%s", Version)

	// Determine directories
	appName := "timezones"
	if *configDir == "" {
		*configDir = os.Getenv("CONFIG_DIR")
		if *configDir == "" {
			*configDir = paths.GetConfigDir(appName)
		}
	}
	if *dataDir == "" {
		*dataDir = os.Getenv("DATA_DIR")
		if *dataDir == "" {
			*dataDir = paths.GetDataDir(appName)
		}
	}
	if *logsDir == "" {
		*logsDir = os.Getenv("LOGS_DIR")
		if *logsDir == "" {
			*logsDir = paths.GetLogsDir(appName)
		}
	}

	// Ensure directories exist
	if err := paths.EnsureDir(*configDir); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}
	if err := paths.EnsureDir(*dataDir); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}
	if err := paths.EnsureDir(filepath.Join(*dataDir, "db")); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}
	if err := paths.EnsureDir(*logsDir); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	log.Printf("Config directory: %s", *configDir)
	log.Printf("Data directory: %s", *dataDir)
	log.Printf("Logs directory: %s", *logsDir)

	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = database.GetDBPath(*dataDir)
	}

	if err := database.Initialize(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Check for admin user
	adminExists, err := database.AdminUserExists()
	if err != nil {
		log.Fatalf("Failed to check for admin user: %v", err)
	}

	// Create admin user if needed
	if !adminExists {
		username := *adminUser
		password := *adminPass

		// Check environment variables
		if username == "" {
			username = os.Getenv("ADMIN_USER")
		}
		if password == "" {
			password = os.Getenv("ADMIN_PASSWORD")
		}

		// Use defaults if still not provided
		if username == "" {
			username = "administrator"
		}
		if password == "" {
			password = "changeme"
			log.Println("⚠️  WARNING: Using default admin password. Please change it!")
		}

		creds, err := database.CreateAdminUser(username, password)
		if err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}

		log.Printf("✓ Admin user created: %s", username)

		// Determine port for credentials file
		portForCreds := *port
		if portForCreds == "" {
			portForCreds = os.Getenv("PORT")
		}
		if portForCreds == "" {
			portForCreds = "8080"
		}

		// Save credentials to file (with password)
		if err := database.SaveCredentialsToFile(creds, *configDir, portForCreds, password); err != nil {
			log.Printf("⚠️  Failed to save credentials file: %v", err)
		} else {
			credPath := filepath.Join(*configDir, "credentials.txt")
			log.Printf("✓ Credentials saved to: %s", credPath)
		}
	}

	// Initialize timezones service with embedded JSON data
	tzService, err := timezones.NewService(timezonesJSON)
	if err != nil {
		log.Fatalf("Failed to initialize timezones service: %v", err)
	}

	log.Printf("✓ Loaded %d timezones", tzService.Count())

	// Determine server configuration
	if *port == "" {
		*port = os.Getenv("PORT")
		if *port == "" {
			*port = "8080"
		}
	}
	if *address == "" {
		*address = os.Getenv("ADDRESS")
		if *address == "" {
			*address = "0.0.0.0"
		}
	}

	// Start HTTP server
	srv := server.New(tzService, *address, *port, Version, BuildDate, Commit)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
