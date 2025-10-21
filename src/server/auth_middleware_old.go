package server

import (
	"net/http"
	"strings"

	"github.com/casjay/timezones/src/database"
)

// authMiddleware validates admin authentication
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		// Check for Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondError(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		token := parts[1]

		// Validate token
		_, err := database.ValidateToken(token)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Token is valid, proceed
		next.ServeHTTP(w, r)
	})
}
