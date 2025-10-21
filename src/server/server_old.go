package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/casjay/timezones/src/timezones"
	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
	router      *mux.Router
	tzService   *timezones.Service
	address     string
	port        string
	version     string
	buildDate   string
	commit      string
}

// New creates a new HTTP server
func New(tzService *timezones.Service, address, port, version, buildDate, commit string) *Server {
	s := &Server{
		router:    mux.NewRouter(),
		tzService: tzService,
		address:   address,
		port:      port,
		version:   version,
		buildDate: buildDate,
		commit:    commit,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Static files
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Public routes
	s.router.HandleFunc("/", s.handleHome).Methods("GET")
	s.router.HandleFunc("/healthz", s.handleHealth).Methods("GET")
	s.router.HandleFunc("/status", s.handleHealth).Methods("GET")

	// API v1 routes
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/timezones.json", s.handleTimezonesJSON).Methods("GET")
	api.HandleFunc("/timezones", s.handleTimezonesAll).Methods("GET")
	api.HandleFunc("/timezones/search", s.handleTimezonesSearch).Methods("GET")
	api.HandleFunc("/timezones/offset/{offset}", s.handleTimezonesByOffset).Methods("GET")
	api.HandleFunc("/timezones/abbr/{abbr}", s.handleTimezonesByAbbr).Methods("GET")
	api.HandleFunc("/timezones/utc/{utc}", s.handleTimezonesByUTC).Methods("GET")
	api.HandleFunc("/timezones/value/{value}", s.handleTimezoneByValue).Methods("GET")
	api.HandleFunc("/stats", s.handleStats).Methods("GET")
	api.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Admin routes (protected by middleware)
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(s.authMiddleware)
	admin.HandleFunc("/settings", s.handleAdminSettings).Methods("GET", "POST")
	admin.HandleFunc("/settings/{key}", s.handleAdminSettingDelete).Methods("DELETE")

	// Middleware
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.corsMiddleware)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.address, s.port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on %s", addr)
	log.Printf("Version: %s (built on %s, commit %s)", s.version, s.buildDate, s.commit)
	log.Printf("Access URL: http://%s:%s", s.address, s.port)

	return srv.ListenAndServe()
}

// loggingMiddleware logs all HTTP requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
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
