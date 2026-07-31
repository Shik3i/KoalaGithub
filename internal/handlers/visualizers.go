package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Shik3i/KoalaGithub/internal/db"
	"github.com/Shik3i/KoalaGithub/internal/models"
)

const maxRequestBody = 1024

var (
	AppVersion = "dev"
	uuidV4     = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	voteLimits = newVoteLimiter(30, time.Hour)
)

type rateEntry struct {
	count   int
	expires time.Time
}

type voteLimiter struct {
	mu        sync.Mutex
	limit     int
	window    time.Duration
	lastSweep time.Time
	entries   map[string]rateEntry
}

func newVoteLimiter(limit int, window time.Duration) *voteLimiter {
	return &voteLimiter{limit: limit, window: window, entries: make(map[string]rateEntry)}
}

func (l *voteLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if now.Sub(l.lastSweep) >= time.Minute {
		for candidate, value := range l.entries {
			if !value.expires.After(now) {
				delete(l.entries, candidate)
			}
		}
		l.lastSweep = now
	}
	if _, exists := l.entries[key]; !exists && len(l.entries) >= 10_000 {
		return false
	}
	entry := l.entries[key]
	if !entry.expires.After(now) {
		entry = rateEntry{expires: now.Add(l.window)}
	}
	if entry.count >= l.limit {
		return false
	}
	entry.count++
	l.entries[key] = entry
	return true
}

func GetVisualizersHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Device-ID")))
	if deviceID != "" && !uuidV4.MatchString(deviceID) {
		writeError(w, http.StatusBadRequest, "invalid device identifier")
		return
	}
	items, err := db.GetAllVisualizers(deviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load visualizers")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func VoteHandler(w http.ResponseWriter, r *http.Request) {
	if !isJSON(r.Header.Get("Content-Type")) {
		writeError(w, http.StatusUnsupportedMediaType, "content type must be application/json")
		return
	}
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" || len(id) > 100 {
		writeError(w, http.StatusBadRequest, "invalid visualizer identifier")
		return
	}
	req, err := decodeVoteRequest(w, r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !voteLimits.allow(clientIP(r) + "|" + req.DeviceID) {
		w.Header().Set("Retry-After", "3600")
		writeError(w, http.StatusTooManyRequests, "too many vote requests")
		return
	}

	voted, count, err := db.ToggleVote(id, req.DeviceID)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "visualizer not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update vote")
		return
	}
	writeJSON(w, http.StatusOK, models.VoteResponse{Voted: voted, VoteCount: count})
}

func DeleteVotesHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Device-ID")))
	if !uuidV4.MatchString(deviceID) {
		writeError(w, http.StatusBadRequest, "invalid device identifier")
		return
	}
	if !voteLimits.allow(clientIP(r) + "|" + deviceID) {
		w.Header().Set("Retry-After", "3600")
		writeError(w, http.StatusTooManyRequests, "too many vote requests")
		return
	}
	deleted, err := db.DeleteDeviceVotes(deviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete votes")
		return
	}
	writeJSON(w, http.StatusOK, models.DeleteVotesResponse{Deleted: deleted})
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status":  "unhealthy",
			"version": AppVersion,
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": AppVersion,
	})
}

func decodeVoteRequest(w http.ResponseWriter, r *http.Request) (models.VoteRequest, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var req models.VoteRequest
	if err := decoder.Decode(&req); err != nil {
		return req, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return req, errors.New("multiple JSON values")
	}
	req.DeviceID = strings.ToLower(strings.TrimSpace(req.DeviceID))
	if !uuidV4.MatchString(req.DeviceID) {
		return req, errors.New("invalid UUID v4")
	}
	return req, nil
}

func isJSON(contentType string) bool {
	contentType = strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	return contentType == "application/json"
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
