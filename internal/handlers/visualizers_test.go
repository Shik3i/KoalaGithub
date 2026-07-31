package handlers

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Shik3i/KoalaGithub/internal/db"
	"github.com/Shik3i/KoalaGithub/internal/models"
)

//go:embed test_visualizers.json
var testSeedJSON []byte

func TestVisualizersAPIAndVoting(t *testing.T) {
	testDB := "./data/test_koalagithub.db"
	_ = os.Remove(testDB)
	defer os.Remove(testDB)

	if err := db.InitDB(testDB, testSeedJSON); err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/visualizers", GetVisualizersHandler)
	mux.HandleFunc("POST /api/visualizers/{id}/vote", VoteHandler)

	// 1. Test GET /api/visualizers
	req := httptest.NewRequest("GET", "/api/visualizers?deviceId=device-123", nil)
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
	voteReqBody, _ := json.Marshal(models.VoteRequest{DeviceID: "device-123"})
	voteReq := httptest.NewRequest("POST", "/api/visualizers/"+targetID+"/vote", bytes.NewBuffer(voteReqBody))
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
	req2 := httptest.NewRequest("GET", "/api/visualizers?deviceId=device-123", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	var visualizers2 []models.Visualizer
	_ = json.Unmarshal(w2.Body.Bytes(), &visualizers2)

	if !visualizers2[0].UserVoted {
		t.Errorf("Expected userVoted=true for device-123")
	}
}
