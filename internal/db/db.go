package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Shik3i/KoalaGithub/internal/models"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

type MemoryCache struct {
	sync.RWMutex
	items     []models.Visualizer
	byID      map[string]*models.Visualizer
	userVotes map[string]map[string]bool
}

var Cache = &MemoryCache{
	byID:      make(map[string]*models.Visualizer),
	userVotes: make(map[string]map[string]bool),
}

func InitDB(dbPath string, seedJSON []byte) error {
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create db directory: %w", err)
		}
	}

	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open sqlite database: %w", err)
	}
	database.SetMaxOpenConns(1)
	if _, err := database.Exec("PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL; PRAGMA busy_timeout = 5000;"); err != nil {
		database.Close()
		return fmt.Errorf("configure sqlite database: %w", err)
	}
	if err := database.Ping(); err != nil {
		database.Close()
		return fmt.Errorf("ping sqlite database: %w", err)
	}
	DB = database
	initialized := false
	defer func() {
		if !initialized {
			_ = database.Close()
			DB = nil
		}
	}()

	if err := createTables(); err != nil {
		return fmt.Errorf("create database tables: %w", err)
	}
	if err := seedInitialData(seedJSON); err != nil {
		return fmt.Errorf("seed visualizers: %w", err)
	}
	if err := LoadCacheFromDB(); err != nil {
		return fmt.Errorf("load RAM cache: %w", err)
	}

	initialized = true
	log.Printf("SQLite database and RAM cache initialized at %s", dbPath)
	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}
	database := DB
	DB = nil
	return database.Close()
}

func Ping() error {
	if DB == nil {
		return fmt.Errorf("database is not initialized")
	}
	return DB.Ping()
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
		image_format TEXT NOT NULL DEFAULT 'svg',
		image_url_template TEXT NOT NULL,
		markdown_template TEXT NOT NULL,
		website_url TEXT NOT NULL,
		repository_url TEXT NOT NULL,
		themes_json TEXT NOT NULL,
		default_theme TEXT NOT NULL,
		requires_username INTEGER NOT NULL,
		requires_external_setup INTEGER NOT NULL,
		setup_explanation TEXT,
		action_workflow_yaml TEXT,
		privacy_notice TEXT,
		added_at TEXT NOT NULL,
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
	CREATE INDEX IF NOT EXISTS idx_votes_device_id ON votes(device_id);
	`
	if _, err := DB.Exec(schema); err != nil {
		return err
	}

	columns := []struct {
		name       string
		definition string
	}{
		{"github_stars", "INTEGER DEFAULT 0"},
		{"last_refreshed_at", "TEXT"},
		{"next_refresh_at", "TEXT"},
		{"image_format", "TEXT NOT NULL DEFAULT 'svg'"},
		{"action_workflow_yaml", "TEXT"},
	}
	for _, column := range columns {
		exists, err := columnExists("visualizers", column.name)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		if _, err := DB.Exec(
			"ALTER TABLE visualizers ADD COLUMN " + column.name + " " + column.definition,
		); err != nil {
			return fmt.Errorf("add visualizers.%s: %w", column.name, err)
		}
	}

	hasLegacyRank, err := columnExists("visualizers", "popularity_rank")
	if err != nil {
		return err
	}
	if hasLegacyRank {
		if _, err := DB.Exec("ALTER TABLE visualizers DROP COLUMN popularity_rank"); err != nil {
			return fmt.Errorf("remove legacy popularity rank: %w", err)
		}
	}
	return nil
}

func columnExists(table, column string) (bool, error) {
	rows, err := DB.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id           int
			name         string
			columnType   string
			notNull      int
			defaultValue sql.NullString
			primaryKey   int
		)
		if err := rows.Scan(&id, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, err
		}
		if name == column {
			return true, nil
		}
	}
	return false, rows.Err()
}

func LoadCacheFromDB() error {
	Cache.Lock()
	defer Cache.Unlock()

	rows, err := DB.Query(`
		SELECT
			v.id, v.name, v.description, v.category, v.tags_json, v.preview_type, v.image_format,
			v.image_url_template, v.markdown_template, v.website_url, v.repository_url,
			v.themes_json, v.default_theme, v.requires_username, v.requires_external_setup,
			COALESCE(v.setup_explanation, ''), COALESCE(v.action_workflow_yaml, ''),
			COALESCE(v.privacy_notice, ''),
			v.added_at, v.enabled, v.estimated_height,
			(SELECT COUNT(*) FROM votes vt WHERE vt.visualizer_id = v.id),
			COALESCE(v.github_stars, 0), COALESCE(v.last_refreshed_at, ''), COALESCE(v.next_refresh_at, '')
		FROM visualizers v
		WHERE v.enabled = 1
		ORDER BY v.name COLLATE NOCASE ASC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	Cache.items = nil
	Cache.byID = make(map[string]*models.Visualizer)
	for rows.Next() {
		var item models.Visualizer
		var tagsJSON, themesJSON string
		var requiresUsername, requiresSetup, enabled int
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &item.Category, &tagsJSON, &item.PreviewType, &item.ImageFormat,
			&item.ImageURLTemplate, &item.MarkdownTemplate, &item.WebsiteURL, &item.RepositoryURL,
			&themesJSON, &item.DefaultTheme, &requiresUsername, &requiresSetup,
			&item.SetupExplanation, &item.ActionWorkflowYaml, &item.PrivacyNotice,
			&item.AddedAt,
			&enabled, &item.EstimatedHeight, &item.VoteCount, &item.GitHubStars,
			&item.LastRefreshedAt, &item.NextRefreshAt,
		); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(tagsJSON), &item.Tags); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(themesJSON), &item.Themes); err != nil {
			return err
		}
		item.RequiresUsername = requiresUsername == 1
		item.RequiresExternalSetup = requiresSetup == 1
		item.Enabled = enabled == 1
		Cache.items = append(Cache.items, item)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for i := range Cache.items {
		Cache.byID[Cache.items[i].ID] = &Cache.items[i]
	}

	voteRows, err := DB.Query("SELECT visualizer_id, device_id FROM votes")
	if err != nil {
		return err
	}
	defer voteRows.Close()
	Cache.userVotes = make(map[string]map[string]bool)
	voteTotal := 0
	for voteRows.Next() {
		var visualizerID, deviceID string
		if err := voteRows.Scan(&visualizerID, &deviceID); err != nil {
			return err
		}
		if Cache.userVotes[deviceID] == nil {
			Cache.userVotes[deviceID] = make(map[string]bool)
		}
		Cache.userVotes[deviceID][visualizerID] = true
		voteTotal++
	}
	if err := voteRows.Err(); err != nil {
		return err
	}

	log.Printf("RAM cache loaded %d visualizers and %d votes", len(Cache.items), voteTotal)
	return nil
}

func GetAllVisualizers(deviceID string) ([]models.Visualizer, error) {
	Cache.RLock()
	defer Cache.RUnlock()

	result := make([]models.Visualizer, len(Cache.items))
	for i, item := range Cache.items {
		result[i] = item
		result[i].UserVoted = Cache.userVotes[deviceID] != nil && Cache.userVotes[deviceID][item.ID]
	}
	return result, nil
}

func ToggleVote(visualizerID, deviceID string) (bool, int, error) {
	Cache.Lock()
	defer Cache.Unlock()

	item, exists := Cache.byID[visualizerID]
	if !exists {
		return false, 0, sql.ErrNoRows
	}

	userVotes := Cache.userVotes[deviceID]
	if userVotes != nil && userVotes[visualizerID] {
		result, err := DB.Exec("DELETE FROM votes WHERE visualizer_id = ? AND device_id = ?", visualizerID, deviceID)
		if err != nil {
			return false, 0, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return false, 0, err
		}
		if affected > 0 {
			delete(userVotes, visualizerID)
			if item.VoteCount > 0 {
				item.VoteCount--
			}
		}
		return false, item.VoteCount, nil
	}

	result, err := DB.Exec("INSERT OR IGNORE INTO votes (visualizer_id, device_id) VALUES (?, ?)", visualizerID, deviceID)
	if err != nil {
		return false, 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, 0, err
	}
	if affected > 0 {
		if Cache.userVotes[deviceID] == nil {
			Cache.userVotes[deviceID] = make(map[string]bool)
		}
		Cache.userVotes[deviceID][visualizerID] = true
		item.VoteCount++
	}
	return true, item.VoteCount, nil
}

func DeleteDeviceVotes(deviceID string) (int, error) {
	Cache.Lock()
	defer Cache.Unlock()

	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT visualizer_id FROM votes WHERE device_id = ?", deviceID)
	if err != nil {
		return 0, err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return 0, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return 0, err
	}
	rows.Close()

	if _, err := tx.Exec("DELETE FROM votes WHERE device_id = ?", deviceID); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	for _, id := range ids {
		if item := Cache.byID[id]; item != nil && item.VoteCount > 0 {
			item.VoteCount--
		}
	}
	delete(Cache.userVotes, deviceID)
	return len(ids), nil
}

func seedInitialData(seedJSON []byte) error {
	if len(seedJSON) == 0 {
		return nil
	}
	var initial []models.Visualizer
	if err := json.Unmarshal(seedJSON, &initial); err != nil {
		return fmt.Errorf("parse visualizers.json: %w", err)
	}

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO visualizers (
			id, name, description, category, tags_json, preview_type, image_format,
			image_url_template, markdown_template, website_url, repository_url,
			themes_json, default_theme, requires_username, requires_external_setup,
			setup_explanation, action_workflow_yaml, privacy_notice, added_at,
			enabled, estimated_height
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name, description = excluded.description, category = excluded.category,
			tags_json = excluded.tags_json, preview_type = excluded.preview_type,
			image_format = excluded.image_format, image_url_template = excluded.image_url_template,
			markdown_template = excluded.markdown_template, website_url = excluded.website_url,
			repository_url = excluded.repository_url, themes_json = excluded.themes_json,
			default_theme = excluded.default_theme, requires_username = excluded.requires_username,
			requires_external_setup = excluded.requires_external_setup,
			setup_explanation = excluded.setup_explanation,
			action_workflow_yaml = excluded.action_workflow_yaml,
			privacy_notice = excluded.privacy_notice,
			added_at = excluded.added_at,
			enabled = excluded.enabled, estimated_height = excluded.estimated_height
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range initial {
		tags, err := json.Marshal(item.Tags)
		if err != nil {
			return err
		}
		themes, err := json.Marshal(item.Themes)
		if err != nil {
			return err
		}
		if _, err := stmt.Exec(
			item.ID, item.Name, item.Description, item.Category, string(tags), item.PreviewType, item.ImageFormat,
			item.ImageURLTemplate, item.MarkdownTemplate, item.WebsiteURL, item.RepositoryURL,
			string(themes), item.DefaultTheme, boolInt(item.RequiresUsername), boolInt(item.RequiresExternalSetup),
			item.SetupExplanation, item.ActionWorkflowYaml, item.PrivacyNotice, item.AddedAt,
			boolInt(item.Enabled), item.EstimatedHeight,
		); err != nil {
			return err
		}
	}

	if len(initial) > 0 {
		placeholders := make([]string, len(initial))
		args := make([]any, len(initial))
		for i, item := range initial {
			placeholders[i] = "?"
			args[i] = item.ID
		}
		if _, err := tx.Exec("DELETE FROM visualizers WHERE id NOT IN ("+strings.Join(placeholders, ",")+")", args...); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func UpdateVisualizerStars(id string, stars int, lastRefreshed, nextRefresh string) error {
	result, err := DB.Exec(`
		UPDATE visualizers
		SET github_stars = ?, last_refreshed_at = ?, next_refresh_at = ?
		WHERE id = ?`,
		stars, lastRefreshed, nextRefresh, id,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	Cache.Lock()
	defer Cache.Unlock()
	if item := Cache.byID[id]; item != nil {
		item.GitHubStars = stars
		item.LastRefreshedAt = lastRefreshed
		item.NextRefreshAt = nextRefresh
	}
	return nil
}
