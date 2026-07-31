package handlers

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/Shik3i/KoalaGithub/internal/db"
	"github.com/Shik3i/KoalaGithub/internal/models"
)

//go:embed test_visualizers.json
var testSeedJSON []byte

func TestVisualizersAPIAndVoting(t *testing.T) {
	const deviceID = "b5bc0a48-04fc-4e0c-bf81-b1fb0cbf63d7"
	testDB := filepath.Join(t.TempDir(), "test_koalagithub.db")

	if err := db.InitDB(testDB, testSeedJSON); err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/visualizers", GetVisualizersHandler)
	mux.HandleFunc("POST /api/visualizers/{id}/vote", VoteHandler)
	mux.HandleFunc("DELETE /api/votes", DeleteVotesHandler)

	// 1. Test GET /api/visualizers
	req := httptest.NewRequest("GET", "/api/visualizers", nil)
	req.Header.Set("X-Device-ID", deviceID)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var visualizers []models.Visualizer
	if err := json.Unmarshal(w.Body.Bytes(), &visualizers); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if len(visualizers) == 0 {
		t.Fatalf("Expected seeded visualizers, got 0")
	}

	targetID := visualizers[0].ID
	initialVotes := visualizers[0].VoteCount

	// 2. Test POST /api/visualizers/{id}/vote (upvote)
	voteReqBody, _ := json.Marshal(models.VoteRequest{DeviceID: deviceID})
	voteReq := httptest.NewRequest("POST", "/api/visualizers/"+targetID+"/vote", bytes.NewBuffer(voteReqBody))
	voteReq.Header.Set("Content-Type", "application/json")
	voteW := httptest.NewRecorder()
	mux.ServeHTTP(voteW, voteReq)

	if voteW.Code != http.StatusOK {
		t.Fatalf("Expected vote status 200, got %d, body: %s", voteW.Code, voteW.Body.String())
	}

	var voteResp models.VoteResponse
	if err := json.Unmarshal(voteW.Body.Bytes(), &voteResp); err != nil {
		t.Fatalf("Failed to parse vote response: %v", err)
	}

	if !voteResp.Voted {
		t.Errorf("Expected voted=true, got false")
	}
	if voteResp.VoteCount != initialVotes+1 {
		t.Errorf("Expected voteCount=%d, got %d", initialVotes+1, voteResp.VoteCount)
	}

	// 3. Test GET /api/visualizers again to verify userVoted=true
	req2 := httptest.NewRequest("GET", "/api/visualizers", nil)
	req2.Header.Set("X-Device-ID", deviceID)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	var visualizers2 []models.Visualizer
	_ = json.Unmarshal(w2.Body.Bytes(), &visualizers2)

	if !visualizers2[0].UserVoted {
		t.Errorf("Expected userVoted=true for current device")
	}

	deleteReq := httptest.NewRequest("DELETE", "/api/votes", nil)
	deleteReq.Header.Set("X-Device-ID", deviceID)
	deleteW := httptest.NewRecorder()
	mux.ServeHTTP(deleteW, deleteReq)
	if deleteW.Code != http.StatusOK {
		t.Fatalf("Expected delete status 200, got %d", deleteW.Code)
	}
}

func TestVoteValidation(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/visualizers/example/vote", bytes.NewBufferString(`{"deviceId":"not-a-uuid"}`))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "example")
	w := httptest.NewRecorder()
	VoteHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}
