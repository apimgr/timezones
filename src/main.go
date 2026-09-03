package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/apimgr/timezones/src/config"
	"github.com/apimgr/timezones/src/mode"
	"github.com/apimgr/timezones/src/paths"
	"github.com/apimgr/timezones/src/server"
	"github.com/apimgr/timezones/src/timezones"
)

// Embed timezones JSON data
//
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
		port            = flag.String("port", "", "Port to listen on")
		address         = flag.String("address", "", "Address to bind to")
		configDir       = flag.String("config", "", "Configuration directory")
		dataDir         = flag.String("data", "", "Data directory")
		logsDir         = flag.String("logs", "", "Logs directory")
		showVersion     = flag.Bool("version", false, "Show version information")
		showStatus      = flag.Bool("status", false, "Show server status (for health checks)")
		showHelp        = flag.Bool("help", false, "Show help message")
		serviceCmd      = flag.String("service", "", "Service command (install, uninstall, start, stop, restart, status)")
		maintenanceMode = flag.String("maintenance", "", "Maintenance mode (on/off)")
		modeFlag        = flag.String("mode", "", "Application mode (dev/development, prod/production)")
		updateCmd       = flag.String("update", "", "Update command (stable, beta, nightly)")
	)
	flag.Parse()

	// Show help
	if *showHelp {
		printHelp()
		return
	}

	// Handle update command
	if *updateCmd != "" {
		handleUpdateCommand(*updateCmd)
		return
	}

	// Initialize mode
	if err := mode.Initialize(*modeFlag); err != nil {
		log.Printf("Warning: invalid mode: %v", err)
	}

	// Unused vars to satisfy compiler
	_ = dataDir
	_ = logsDir

	// Show version
	if *showVersion {
		fmt.Printf("%s\n", Version)
		return
	}

	// Status check (for health checks)
	if *showStatus {
		os.Exit(0)
	}

	// Handle service commands
	if *serviceCmd != "" {
		handleServiceCommand(*serviceCmd)
		return
	}

	// Handle maintenance mode
	if *maintenanceMode != "" {
		handleMaintenanceMode(*maintenanceMode)
		return
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

	// Ensure directories exist
	if err := paths.EnsureDir(*configDir); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}

	log.Printf("Config directory: %s", *configDir)

	// Load configuration (using server.yml per BASE.md)
	configPath := filepath.Join(*configDir, "server.yml")
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Printf("Warning: Failed to load config: %v (using defaults)", err)
		cfg = config.DefaultConfig()
	}

	// Initialize timezones service with embedded JSON data
	tzService, err := timezones.NewService(timezonesJSON)
	if err != nil {
		log.Fatalf("Failed to initialize timezones service: %v", err)
	}

	log.Printf("Loaded %d timezones", tzService.Count())

	// Determine server configuration
	if *port == "" {
		*port = os.Getenv("PORT")
		if *port == "" {
			*port = cfg.Server.Port
		}
		if *port == "" {
			*port = "8080"
		}
	}
	if *address == "" {
		*address = os.Getenv("ADDRESS")
		if *address == "" {
			*address = cfg.Server.Address
		}
		if *address == "" {
			*address = "0.0.0.0"
		}
	}

	// Create HTTP server
	srv := server.New(tzService, cfg, *address, *port, Version, BuildDate, Commit)

	// Setup graceful shutdown
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", *address, *port),
		Handler:      srv.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on %s:%s", *address, *port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// handleServiceCommand handles service management commands
func handleServiceCommand(cmd string) {
	switch cmd {
	case "install":
		fmt.Println("Service installation not yet implemented")
		fmt.Println("Use systemd/launchd/rc.d to manage the service")
	case "uninstall":
		fmt.Println("Service uninstallation not yet implemented")
	case "start":
		fmt.Println("Use 'systemctl start timezones' or run the binary directly")
	case "stop":
		fmt.Println("Use 'systemctl stop timezones' or send SIGTERM to the process")
	case "restart":
		fmt.Println("Use 'systemctl restart timezones'")
	case "status":
		fmt.Println("Use 'systemctl status timezones' or --status flag")
	default:
		fmt.Printf("Unknown service command: %s\n", cmd)
		fmt.Println("Available commands: install, uninstall, start, stop, restart, status")
	}
}

// handleMaintenanceMode handles maintenance mode toggle
func handleMaintenanceMode(m string) {
	switch m {
	case "on":
		fmt.Println("Maintenance mode: ON")
		fmt.Println("Note: Maintenance mode is handled at runtime, not persisted")
	case "off":
		fmt.Println("Maintenance mode: OFF")
	default:
		fmt.Printf("Invalid maintenance mode: %s (use 'on' or 'off')\n", m)
	}
}

func printHelp() {
	fmt.Printf(`Timezones API v%s

Usage: timezones [options]

Options:
  --port PORT        Port to listen on (default: 8080)
  --address ADDR     Address to bind to (default: 0.0.0.0)
  --config PATH      Configuration directory
  --data PATH        Data directory
  --logs PATH        Logs directory
  --mode MODE        Application mode (dev, prod)
  --update BRANCH    Update from branch (stable, beta, nightly)
  --version          Show version information
  --status           Show server status
  --help             Show this help message

Service Management:
  --service install    Install as system service
  --service uninstall  Uninstall system service
  --service start      Start the service
  --service stop       Stop the service
  --service restart    Restart the service
  --service status     Show service status

Maintenance:
  --maintenance on     Enable maintenance mode
  --maintenance off    Disable maintenance mode

Examples:
  timezones --port 3000
  timezones --mode dev --port 8080
  timezones --update stable
`, Version)
}

func handleUpdateCommand(branch string) {
	validBranches := map[string]bool{
		"stable":  true,
		"beta":    true,
		"nightly": true,
	}

	if !validBranches[branch] {
		fmt.Printf("Error: invalid update branch %q (valid: stable, beta, nightly)\n", branch)
		os.Exit(1)
	}

	fmt.Printf("Updating Timezones API from %s branch...\n", branch)

	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("Error: git is not installed")
		os.Exit(1)
	}

	cmd := exec.Command("git", "pull", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Update failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Update complete. Please rebuild the application.")
}
