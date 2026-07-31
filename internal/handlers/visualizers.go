package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Shik3i/KoalaGithub/internal/db"
	"github.com/Shik3i/KoalaGithub/internal/models"
)

func GetVisualizersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deviceID := r.Header.Get("X-Device-ID")
	if deviceID == "" {
		deviceID = r.URL.Query().Get("deviceId")
	}

	items, err := db.GetAllVisualizers(deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func VoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract visualizer ID from path e.g. /api/visualizers/{id}/vote
	id := r.PathValue("id")
	if id == "" {
		p := r.URL.Path
		if strings.HasPrefix(p, "/api/visualizers/") {
			p = strings.TrimPrefix(p, "/api/visualizers/")
			if idx := strings.Index(p, "/vote"); idx != -1 {
				id = p[:idx]
			}
		}
	}
	if id == "" {
		http.Error(w, "Missing visualizer ID", http.StatusBadRequest)
		return
	}

	var req models.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DeviceID == "" {
		http.Error(w, "Invalid request body; deviceId required", http.StatusBadRequest)
		return
	}

	voted, newCount, err := db.ToggleVote(id, req.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := models.VoteResponse{
		Voted:     voted,
		VoteCount: newCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
