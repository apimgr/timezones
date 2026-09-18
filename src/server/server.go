package server

import (
	"net"
	"net/http"
	"strings"
	"time"

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

	// Real IP (chi's middleware.RealIP is deprecated: it trusts forwarded
	// headers unconditionally, which allows IP spoofing. Only trust them
	// when the immediate TCP peer is a private/loopback address, per
	// AI.md PART 12 "Trusted Proxies".)
	s.router.Use(realIPMiddleware)

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
		// Wildcard "*" rather than {utc}: IANA identifiers like
		// "America/New_York" contain slashes, which decode into extra path
		// segments a single-segment chi param can never match.
		r.Get("/timezones/utc/*", s.handleTimezonesByUTC)
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

// realIPMiddleware rewrites r.RemoteAddr with the client IP resolved from
// X-Forwarded-For/X-Real-IP, but only when the immediate TCP peer is a
// trusted (private/loopback) address. Unlike chi's deprecated
// middleware.RealIP, forwarded headers from an untrusted peer are ignored,
// preventing IP spoofing (AI.md PART 12 "Trusted Proxies").
func realIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isTrustedPeer(r.RemoteAddr) {
			if ip := clientIPFromHeaders(r); ip != "" {
				r.RemoteAddr = ip
			}
		}
		next.ServeHTTP(w, r)
	})
}

// clientIPFromHeaders extracts the client IP from trusted proxy headers.
func clientIPFromHeaders(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}
	return ""
}

// isTrustedPeer reports whether addr (host:port or bare host) is a
// loopback or private-range address.
func isTrustedPeer(addr string) bool {
	host := addr
	if h, _, err := net.SplitHostPort(addr); err == nil {
		host = h
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified()
}
