package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/casjay/timezones/src/security"
	"github.com/casjay/timezones/src/timezones"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server represents the HTTP server
type Server struct {
	router    *chi.Mux
	tzService *timezones.Service
	address   string
	port      string
	version   string
	buildDate string
	commit    string
}

// New creates a new HTTP server
func New(tzService *timezones.Service, address, port, version, buildDate, commit string) *Server {
	s := &Server{
		router:    chi.NewRouter(),
		tzService: tzService,
		address:   address,
		port:      port,
		version:   version,
		buildDate: buildDate,
		commit:    commit,
	}

	s.setupMiddleware()
	s.setupRoutes()
	return s
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

	// Global rate limiting
	s.router.Use(security.GlobalRateLimiter())

	// Security headers
	s.router.Use(security.SecurityHeadersMiddleware)

	// CORS
	s.router.Use(s.corsMiddleware)
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Static files
	fileServer := http.FileServer(http.FS(staticFS))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Public routes
	s.router.Get("/", s.handleHome)
	s.router.Get("/healthz", s.handleHealth)
	s.router.Get("/status", s.handleHealth)

	// API v1 routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Apply API rate limiting
		r.Use(security.APIRateLimiter())

		// Timezone endpoints
		r.Get("/timezones.json", s.handleTimezonesJSON)
		r.Get("/timezones", s.handleTimezonesAll)
		r.Get("/timezones/search", s.handleTimezonesSearch)
		r.Get("/timezones/offset/{offset}", s.handleTimezonesByOffset)
		r.Get("/timezones/abbr/{abbr}", s.handleTimezonesByAbbr)
		r.Get("/timezones/utc/{utc}", s.handleTimezonesByUTC)
		r.Get("/timezones/value/{value}", s.handleTimezoneByValue)

		// Stats & health
		r.Get("/stats", s.handleStats)
		r.Get("/health", s.handleHealth)

		// Admin routes (protected)
		r.Group(func(r chi.Router) {
			r.Use(security.AdminRateLimiter())
			r.Use(s.authMiddleware)

			r.Get("/admin/settings", s.handleAdminSettings)
			r.Post("/admin/settings", s.handleAdminSettings)
			r.Delete("/admin/settings/{key}", s.handleAdminSettingDelete)
		})
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.address, s.port)

	srv := &http.Server{
		Addr:           addr,
		Handler:        s.router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	log.Printf("Starting server on %s", addr)
	log.Printf("Version: %s (built on %s, commit %s)", s.version, s.buildDate, s.commit)

	return srv.ListenAndServe()
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
