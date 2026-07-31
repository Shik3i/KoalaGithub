package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/Shik3i/KoalaGithub/internal/models"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

type MemoryCache struct {
	sync.RWMutex
	items     []models.Visualizer
	byID      map[string]*models.Visualizer
	userVotes map[string]map[string]bool // deviceID -> map[visualizerID]bool
}

var Cache = &MemoryCache{
	byID:      make(map[string]*models.Visualizer),
	userVotes: make(map[string]map[string]bool),
}

func InitDB(dbPath string, seedJSON []byte) error {
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create db directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	DB = db

	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create database tables: %w", err)
	}

	if err := seedInitialData(seedJSON); err != nil {
		log.Printf("Warning: error seeding initial data: %v", err)
	}

	// Load complete dataset into RAM cache
	if err := LoadCacheFromDB(); err != nil {
		return fmt.Errorf("failed to load RAM cache: %w", err)
	}

	log.Println("SQLite database and RAM Cache initialized successfully at", dbPath)
	return nil
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS visualizers (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		category TEXT NOT NULL,
		tags_json TEXT NOT NULL,
		preview_type TEXT NOT NULL,
		image_url_template TEXT NOT NULL,
		markdown_template TEXT NOT NULL,
		website_url TEXT NOT NULL,
		repository_url TEXT NOT NULL,
		themes_json TEXT NOT NULL,
		default_theme TEXT NOT NULL,
		requires_username INTEGER NOT NULL,
		requires_external_setup INTEGER NOT NULL,
		setup_explanation TEXT,
		privacy_notice TEXT,
		added_at TEXT NOT NULL,
		popularity_rank INTEGER NOT NULL,
		enabled INTEGER NOT NULL,
		estimated_height INTEGER NOT NULL,
		github_stars INTEGER DEFAULT 0,
		last_refreshed_at TEXT,
		next_refresh_at TEXT
	);

	CREATE TABLE IF NOT EXISTS votes (
		visualizer_id TEXT NOT NULL,
		device_id TEXT NOT NULL,
		voted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (visualizer_id, device_id),
		FOREIGN KEY (visualizer_id) REFERENCES visualizers(id) ON DELETE CASCADE
	);
	`
	_, err := DB.Exec(schema)
	if err != nil {
		return err
	}

	// Migrations for existing DBs
	_, _ = DB.Exec("ALTER TABLE visualizers ADD COLUMN github_stars INTEGER DEFAULT 0;")
	_, _ = DB.Exec("ALTER TABLE visualizers ADD COLUMN last_refreshed_at TEXT;")
	_, _ = DB.Exec("ALTER TABLE visualizers ADD COLUMN next_refresh_at TEXT;")

	return nil
}

func LoadCacheFromDB() error {
	Cache.Lock()
	defer Cache.Unlock()

	// 1. Fetch visualizers & total votes from SQLite
	query := `
		SELECT 
			v.id, v.name, v.description, v.category, v.tags_json, v.preview_type,
			v.image_url_template, v.markdown_template, v.website_url, v.repository_url,
			v.themes_json, v.default_theme, v.requires_username, v.requires_external_setup,
			COALESCE(v.setup_explanation, ''), COALESCE(v.privacy_notice, ''),
			v.added_at, v.popularity_rank, v.enabled, v.estimated_height,
			(SELECT COUNT(*) FROM votes vt WHERE vt.visualizer_id = v.id) AS vote_count,
			COALESCE(v.github_stars, 0), COALESCE(v.last_refreshed_at, ''), COALESCE(v.next_refresh_at, '')
		FROM visualizers v
		WHERE v.enabled = 1
		ORDER BY v.popularity_rank ASC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	Cache.items = nil
	Cache.byID = make(map[string]*models.Visualizer)

	for rows.Next() {
		var item models.Visualizer
		var tagsJSON, themesJSON string
		var reqUser, reqSetup, enabled int

		err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &item.Category, &tagsJSON, &item.PreviewType,
			&item.ImageURLTemplate, &item.MarkdownTemplate, &item.WebsiteURL, &item.RepositoryURL,
			&themesJSON, &item.DefaultTheme, &reqUser, &reqSetup,
			&item.SetupExplanation, &item.PrivacyNotice,
			&item.AddedAt, &item.PopularityRank, &enabled, &item.EstimatedHeight,
			&item.VoteCount, &item.GitHubStars, &item.LastRefreshedAt, &item.NextRefreshAt,
		)
		if err != nil {
			return err
		}

		_ = json.Unmarshal([]byte(tagsJSON), &item.Tags)
		_ = json.Unmarshal([]byte(themesJSON), &item.Themes)
		item.RequiresUsername = reqUser == 1
		item.RequiresExternalSetup = reqSetup == 1
		item.Enabled = enabled == 1

		Cache.items = append(Cache.items, item)
	}

	for i := range Cache.items {
		Cache.byID[Cache.items[i].ID] = &Cache.items[i]
	}

	// 2. Load votes mapping into RAM
	voteRows, err := DB.Query("SELECT visualizer_id, device_id FROM votes")
	if err != nil {
		return err
	}
	defer voteRows.Close()

	Cache.userVotes = make(map[string]map[string]bool)
	voteTotal := 0
	for voteRows.Next() {
		var vID, dID string
		if err := voteRows.Scan(&vID, &dID); err == nil {
			if Cache.userVotes[dID] == nil {
				Cache.userVotes[dID] = make(map[string]bool)
			}
			Cache.userVotes[dID][vID] = true
			voteTotal++
		}
	}

	log.Printf("RAM Cache loaded %d visualizers and %d votes from DB (0 disk I/O on reads)", len(Cache.items), voteTotal)
	return nil
}

// GetAllVisualizers returns visualizers directly from RAM cache with 0 disk I/O
func GetAllVisualizers(deviceID string) ([]models.Visualizer, error) {
	Cache.RLock()
	defer Cache.RUnlock()

	result := make([]models.Visualizer, len(Cache.items))
	userVoteMap := Cache.userVotes[deviceID]

	for i, item := range Cache.items {
		result[i] = item
		if userVoteMap != nil && userVoteMap[item.ID] {
			result[i].UserVoted = true
		} else {
			result[i].UserVoted = false
		}
	}

	return result, nil
}

// ToggleVote updates RAM cache immediately and syncs SQLite asynchronously/synchronously
func ToggleVote(visualizerID, deviceID string) (bool, int, error) {
	Cache.Lock()
	defer Cache.Unlock()

	item, exists := Cache.byID[visualizerID]
	if !exists {
		return false, 0, fmt.Errorf("visualizer %s not found", visualizerID)
	}

	if Cache.userVotes[deviceID] == nil {
		Cache.userVotes[deviceID] = make(map[string]bool)
	}

	voted := false
	if Cache.userVotes[deviceID][visualizerID] {
		// Un-vote in RAM & DB
		delete(Cache.userVotes[deviceID], visualizerID)
		if item.VoteCount > 0 {
			item.VoteCount--
		}
		voted = false

		_, err := DB.Exec("DELETE FROM votes WHERE visualizer_id = ? AND device_id = ?", visualizerID, deviceID)
		if err != nil {
			return false, 0, err
		}
	} else {
		// Vote in RAM & DB
		Cache.userVotes[deviceID][visualizerID] = true
		item.VoteCount++
		voted = true

		_, err := DB.Exec("INSERT OR IGNORE INTO votes (visualizer_id, device_id) VALUES (?, ?)", visualizerID, deviceID)
		if err != nil {
			return false, 0, err
		}
	}

	return voted, item.VoteCount, nil
}

func seedInitialData(seedJSON []byte) error {
	if len(seedJSON) == 0 {
		return nil
	}

	var initialVisualizers []models.Visualizer
	if err := json.Unmarshal(seedJSON, &initialVisualizers); err != nil {
		return fmt.Errorf("failed to parse visualizers.json: %w", err)
	}

	// Sync visualizers table with central JSON dataset
	_, err := DB.Exec("DELETE FROM visualizers")
	if err != nil {
		return err
	}

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO visualizers (
			id, name, description, category, tags_json, preview_type,
			image_url_template, markdown_template, website_url, repository_url,
			themes_json, default_theme, requires_username, requires_external_setup,
			setup_explanation, privacy_notice, added_at, popularity_rank, enabled, estimated_height
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range initialVisualizers {
		tagsBytes, _ := json.Marshal(item.Tags)
		themesBytes, _ := json.Marshal(item.Themes)
		reqUser := 0
		if item.RequiresUsername {
			reqUser = 1
		}
		reqSetup := 0
		if item.RequiresExternalSetup {
			reqSetup = 1
		}
		enabled := 0
		if item.Enabled {
			enabled = 1
		}

		_, err = stmt.Exec(
			item.ID, item.Name, item.Description, item.Category, string(tagsBytes), item.PreviewType,
			item.ImageURLTemplate, item.MarkdownTemplate, item.WebsiteURL, item.RepositoryURL,
			string(themesBytes), item.DefaultTheme, reqUser, reqSetup,
			item.SetupExplanation, item.PrivacyNotice, item.AddedAt, item.PopularityRank, enabled, item.EstimatedHeight,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// UpdateVisualizerStars updates RAM cache and persists star count + refresh timestamps to SQLite
func UpdateVisualizerStars(id string, stars int, lastRefreshed, nextRefresh string) error {
	Cache.Lock()
	if item, ok := Cache.byID[id]; ok {
		item.GitHubStars = stars
		item.LastRefreshedAt = lastRefreshed
		item.NextRefreshAt = nextRefresh
	}
	Cache.Unlock()

	_, err := DB.Exec(`
		UPDATE visualizers 
		SET github_stars = ?, last_refreshed_at = ?, next_refresh_at = ? 
		WHERE id = ?`,
		stars, lastRefreshed, nextRefresh, id,
	)
	return err
}
