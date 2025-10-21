package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
			"status":      "healthy",
			"version":     s.version,
			"timezones":   s.tzService.Count(),
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

// handleTimezonesByOffset returns timezones by UTC offset
func (s *Server) handleTimezonesByOffset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	offsetStr := vars["offset"]

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
	vars := mux.Vars(r)
	abbr := vars["abbr"]

	results := s.tzService.GetByAbbr(abbr)
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    results,
	})
}

// handleTimezonesByUTC returns timezone by UTC identifier
func (s *Server) handleTimezonesByUTC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	utc := vars["utc"]

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
	vars := mux.Vars(r)
	value := vars["value"]

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

// handleStats returns statistics about timezone data
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    s.tzService.GetStats(),
	})
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
