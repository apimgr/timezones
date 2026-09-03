package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apimgr/timezones/src/config"
	"github.com/go-chi/chi/v5"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// APIResponse represents a standardized API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// handleHome serves the homepage
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":      "Timezones API",
		"Version":    s.version,
		"BuildDate":  s.buildDate,
		"Commit":     s.commit,
		"TotalCount": s.tzService.Count(),
		"Theme":      config.GetTheme(),
	}

	if err := renderTemplate(w, "home.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// handleHealth returns server health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"status":    "healthy",
			"version":   s.version,
			"timezones": s.tzService.Count(),
		},
	})
}

// handleTimezonesJSON returns the raw JSON file
func (s *Server) handleTimezonesJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(s.tzService.GetRawJSON())
}

// handleTimezonesAll returns all timezones
func (s *Server) handleTimezonesAll(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    s.tzService.GetAll(),
	})
}

// handleTimezonesAllTxt returns all timezones as plain text
func (s *Server) handleTimezonesAllTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	timezones := s.tzService.GetAll()
	for _, tz := range timezones {
		fmt.Fprintf(w, "%s (%s) UTC%+.0f\n", tz.Value, tz.Abbr, tz.Offset)
	}
}

// handleTimezonesSearch searches timezones
func (s *Server) handleTimezonesSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondError(w, http.StatusBadRequest, "Missing search query parameter 'q'")
		return
	}

	results := s.tzService.Search(query)
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    results,
	})
}

// handleTimezonesSearchTxt searches timezones and returns plain text
func (s *Server) handleTimezonesSearchTxt(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error: Missing search query parameter 'q'")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	results := s.tzService.Search(query)
	for _, tz := range results {
		fmt.Fprintf(w, "%s (%s) UTC%+.0f\n", tz.Value, tz.Abbr, tz.Offset)
	}
}

// handleTimezonesByOffset returns timezones by UTC offset
func (s *Server) handleTimezonesByOffset(w http.ResponseWriter, r *http.Request) {
	offsetStr := chi.URLParam(r, "offset")

	offset, err := strconv.ParseFloat(offsetStr, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid offset value")
		return
	}

	results := s.tzService.GetByOffset(offset)
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    results,
	})
}

// handleTimezonesByAbbr returns timezones by abbreviation
func (s *Server) handleTimezonesByAbbr(w http.ResponseWriter, r *http.Request) {
	abbr := chi.URLParam(r, "abbr")

	results := s.tzService.GetByAbbr(abbr)
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    results,
	})
}

// handleTimezonesByUTC returns timezone by UTC identifier
func (s *Server) handleTimezonesByUTC(w http.ResponseWriter, r *http.Request) {
	utc := chi.URLParam(r, "utc")

	result := s.tzService.GetByUTC(utc)
	if result == nil {
		respondError(w, http.StatusNotFound, "Timezone not found")
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    result,
	})
}

// handleTimezoneByValue returns timezone by value
func (s *Server) handleTimezoneByValue(w http.ResponseWriter, r *http.Request) {
	value := chi.URLParam(r, "value")

	result := s.tzService.GetByValue(value)
	if result == nil {
		respondError(w, http.StatusNotFound, "Timezone not found")
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    result,
	})
}

// handleTimezonesRandom returns a random timezone
func (s *Server) handleTimezonesRandom(w http.ResponseWriter, r *http.Request) {
	timezones := s.tzService.GetAll()
	if len(timezones) == 0 {
		respondError(w, http.StatusNotFound, "No timezones available")
		return
	}

	randomTz := timezones[rand.Intn(len(timezones))]
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    randomTz,
	})
}

// handleTimezonesRandomTxt returns a random timezone as plain text
func (s *Server) handleTimezonesRandomTxt(w http.ResponseWriter, r *http.Request) {
	timezones := s.tzService.GetAll()
	if len(timezones) == 0 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "No timezones available")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	randomTz := timezones[rand.Intn(len(timezones))]
	fmt.Fprintf(w, "%s (%s) UTC%+.0f\n", randomTz.Value, randomTz.Abbr, randomTz.Offset)
	fmt.Fprintf(w, "Text: %s\n", randomTz.Text)
	if len(randomTz.UTC) > 0 {
		fmt.Fprintf(w, "UTC Identifiers: %s\n", strings.Join(randomTz.UTC, ", "))
	}
}

// handleStats returns statistics about timezone data
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := s.tzService.GetStats()
	stats["version"] = s.version
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    stats,
	})
}

// handleStatsTxt returns statistics as plain text
func (s *Server) handleStatsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	stats := s.tzService.GetStats()
	fmt.Fprintf(w, "Timezones API Statistics\n")
	fmt.Fprintf(w, "========================\n")
	fmt.Fprintf(w, "Version: %s\n", s.version)
	fmt.Fprintf(w, "Total Timezones: %v\n", stats["total_timezones"])
	fmt.Fprintf(w, "Total UTC Entries: %v\n", stats["total_utc_entries"])
	fmt.Fprintf(w, "DST Timezones: %v\n", stats["dst_timezones"])
	fmt.Fprintf(w, "Non-DST Timezones: %v\n", stats["non_dst_timezones"])
}

// handleCount returns the total count
func (s *Server) handleCount(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"count": s.tzService.Count(),
		},
	})
}

// handleCountTxt returns the count as plain text
func (s *Server) handleCountTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%d\n", s.tzService.Count())
}

// handleManifest serves the PWA manifest
func (s *Server) handleManifest(w http.ResponseWriter, r *http.Request) {
	manifest := map[string]interface{}{
		"name":             "Timezones API",
		"short_name":       "Timezones",
		"description":      "Timezone lookup API",
		"start_url":        "/",
		"display":          "standalone",
		"background_color": "#1a1a1a",
		"theme_color":      "#0066cc",
		"icons": []map[string]string{
			{"src": "/static/images/icon-192.png", "sizes": "192x192", "type": "image/png"},
			{"src": "/static/images/icon-512.png", "sizes": "512x512", "type": "image/png"},
		},
	}
	respondJSON(w, http.StatusOK, manifest)
}

// handleServiceWorker serves the service worker
func (s *Server) handleServiceWorker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte(`// Service Worker for Timezones API
const CACHE_NAME = 'timezones-api-v1';
const urlsToCache = [
  '/',
  '/static/css/main.css',
  '/static/js/main.js'
];

self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(cache => cache.addAll(urlsToCache))
  );
});

self.addEventListener('fetch', event => {
  event.respondWith(
    caches.match(event.request)
      .then(response => response || fetch(event.request))
  );
});
`))
}

// handleRobotsTxt serves robots.txt
func (s *Server) handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "User-agent: *")
	for _, path := range cfg.WebRobots.Allow {
		fmt.Fprintf(w, "Allow: %s\n", path)
	}
	for _, path := range cfg.WebRobots.Deny {
		fmt.Fprintf(w, "Disallow: %s\n", path)
	}
}

// handleSecurityTxt serves security.txt
func (s *Server) handleSecurityTxt(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "# Security Policy")
	if cfg.WebSecurity.Admin != "" {
		fmt.Fprintf(w, "Contact: mailto:%s\n", cfg.WebSecurity.Admin)
	} else {
		fmt.Fprintln(w, "Contact: mailto:security@example.com")
	}
	fmt.Fprintln(w, "Preferred-Languages: en")
}

// respondJSON writes a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError writes an error JSON response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    status,
			Message: message,
		},
	})
}
