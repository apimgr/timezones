package server

import (
	"net/http"
	"time"

	"github.com/apimgr/timezones/src/admin"
	"github.com/apimgr/timezones/src/config"
	"github.com/apimgr/timezones/src/timezones"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server represents the HTTP server
type Server struct {
	router       *chi.Mux
	tzService    *timezones.Service
	config       *config.Config
	adminHandler *admin.Handler
	address      string
	port         string
	version      string
	buildDate    string
	commit       string
}

// New creates a new HTTP server
func New(tzService *timezones.Service, cfg *config.Config, address, port, version, buildDate, commit string) *Server {
	s := &Server{
		router:    chi.NewRouter(),
		tzService: tzService,
		config:    cfg,
		address:   address,
		port:      port,
		version:   version,
		buildDate: buildDate,
		commit:    commit,
	}

	// Initialize admin handler
	s.adminHandler = admin.NewHandler(
		cfg.Server.Admin.Username,
		cfg.Server.Admin.Password,
		cfg.Server.Admin.APIToken,
		cfg.Server.Session.Timeout,
		false, // SSL enabled
		version,
		commit,
		buildDate,
	)

	s.setupMiddleware()
	s.setupRoutes()
	return s
}

// Router returns the chi router
func (s *Server) Router() *chi.Mux {
	return s.router
}

// setupMiddleware configures global middleware
func (s *Server) setupMiddleware() {
	// Recovery from panics
	s.router.Use(middleware.Recoverer)

	// Request ID
	s.router.Use(middleware.RequestID)

	// Real IP
	s.router.Use(middleware.RealIP)

	// Logger
	s.router.Use(middleware.Logger)

	// Timeout
	s.router.Use(middleware.Timeout(60 * time.Second))

	// Throttle concurrent requests
	s.router.Use(middleware.Throttle(1000))

	// CORS
	s.router.Use(s.corsMiddleware)

	// Security headers
	s.router.Use(s.securityHeadersMiddleware)
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Register admin routes
	s.adminHandler.RegisterRoutes(s.router)

	// Static files
	fileServer := http.FileServer(http.FS(staticFS))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Public routes
	s.router.Get("/", s.handleHome)
	s.router.Get("/healthz", s.handleHealth)
	s.router.Get("/health", s.handleHealth)
	s.router.Get("/status", s.handleHealth)

	// PWA support
	s.router.Get("/manifest.json", s.handleManifest)
	s.router.Get("/sw.js", s.handleServiceWorker)
	s.router.Get("/robots.txt", s.handleRobotsTxt)
	s.router.Get("/security.txt", s.handleSecurityTxt)
	s.router.Get("/.well-known/security.txt", s.handleSecurityTxt)

	// API v1 routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Timezone endpoints - JSON
		r.Get("/timezones.json", s.handleTimezonesJSON)
		r.Get("/timezones", s.handleTimezonesAll)
		r.Get("/timezones/search", s.handleTimezonesSearch)
		r.Get("/timezones/offset/{offset}", s.handleTimezonesByOffset)
		r.Get("/timezones/abbr/{abbr}", s.handleTimezonesByAbbr)
		r.Get("/timezones/utc/{utc}", s.handleTimezonesByUTC)
		r.Get("/timezones/value/{value}", s.handleTimezoneByValue)
		r.Get("/timezones/random", s.handleTimezonesRandom)

		// Timezone endpoints - Plain text (.txt)
		r.Get("/timezones.txt", s.handleTimezonesAllTxt)
		r.Get("/timezones/search.txt", s.handleTimezonesSearchTxt)
		r.Get("/timezones/random.txt", s.handleTimezonesRandomTxt)

		// Stats & health
		r.Get("/stats", s.handleStats)
		r.Get("/stats.txt", s.handleStatsTxt)
		r.Get("/health", s.handleHealth)
		r.Get("/count", s.handleCount)
		r.Get("/count.txt", s.handleCountTxt)
	})

	// Shorthand routes
	s.router.Get("/random", s.handleTimezonesRandom)
	s.router.Get("/random.txt", s.handleTimezonesRandomTxt)
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors := config.GetCORS()
		w.Header().Set("Access-Control-Allow-Origin", cors)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// securityHeadersMiddleware adds security headers
func (s *Server) securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
