package server

import (
	"encoding/json"
	"net/http"

	"github.com/casjay/timezones/src/database"
	"github.com/gorilla/mux"
)

// handleAdminSettings handles GET (retrieve all settings) and POST (update setting)
func (s *Server) handleAdminSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		settings, err := database.GetAllSettings()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to retrieve settings")
			return
		}

		respondJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data:    settings,
		})

	case "POST":
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Key == "" || req.Value == "" {
			respondError(w, http.StatusBadRequest, "Missing key or value")
			return
		}

		if err := database.SetSetting(req.Key, req.Value); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to save setting")
			return
		}

		respondJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data: map[string]string{
				"message": "Setting saved successfully",
			},
		})
	}
}

// handleAdminSettingDelete handles DELETE (remove a setting)
func (s *Server) handleAdminSettingDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		respondError(w, http.StatusBadRequest, "Missing key")
		return
	}

	if err := database.DeleteSetting(key); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete setting")
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]string{
			"message": "Setting deleted successfully",
		},
	})
}
