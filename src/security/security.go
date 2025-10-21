package security

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/httprate"
)

// Rate limiting configuration
const (
	GlobalRPS   = 100 // 100 requests per second
	GlobalBurst = 200
	APIRPS      = 50
	APIBurst    = 100
	AdminRPS    = 10
	AdminBurst  = 20
)

// FailedAttempts tracks failed login attempts for brute force protection
var FailedAttempts = sync.Map{}

// SetupRateLimiting returns rate limiting middleware
func GlobalRateLimiter() func(http.Handler) http.Handler {
	return httprate.LimitByIP(GlobalRPS, time.Second)
}

func APIRateLimiter() func(http.Handler) http.Handler {
	return httprate.LimitByIP(APIRPS, time.Second)
}

func AdminRateLimiter() func(http.Handler) http.Handler {
	return httprate.LimitByIP(AdminRPS, time.Second)
}

// SecurityHeadersMiddleware adds security headers to all responses
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// XSS Protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Content Security Policy
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline' https://unpkg.com; "+
				"style-src 'self' 'unsafe-inline' https://unpkg.com; "+
				"img-src 'self' data:; "+
				"font-src 'self' https://unpkg.com; "+
				"connect-src 'self'")

		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// HSTS (if using HTTPS)
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		next.ServeHTTP(w, r)
	})
}

// CheckBruteForce checks if IP has too many failed attempts
func CheckBruteForce(ip string) bool {
	attempts, _ := FailedAttempts.LoadOrStore(ip, 0)
	return attempts.(int) >= 5
}

// RecordFailedAttempt records a failed login attempt
func RecordFailedAttempt(ip string) {
	attempts, _ := FailedAttempts.LoadOrStore(ip, 0)
	FailedAttempts.Store(ip, attempts.(int)+1)
}

// ResetFailedAttempts resets failed attempts for an IP
func ResetFailedAttempts(ip string) {
	FailedAttempts.Delete(ip)
}

// GetClientIP extracts the real client IP from request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (proxy/load balancer)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	return r.RemoteAddr
}

// ValidatePassword validates password complexity
func ValidatePassword(password string) error {
	if len(password) < 12 {
		return http.ErrAbortHandler // Simple error for now
	}
	return nil
}
