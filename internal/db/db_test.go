package db

import (
	"path/filepath"
	"testing"
)

const testSeed = `[
	{
		"id":"one","name":"One","description":"Test","category":"stats","tags":["test"],
		"previewType":"image","imageFormat":"svg","imageUrlTemplate":"https://example.com/{username}.svg",
		"markdownTemplate":"![test](https://example.com/{username}.svg)","websiteUrl":"https://example.com",
		"repositoryUrl":"https://github.com/example/one","themes":[],"defaultTheme":"",
		"requiresUsername":true,"requiresExternalSetup":false,"addedAt":"2026-07-31",
		"enabled":true,"estimatedHeight":100
	}
]`

func TestReseedPreservesVotes(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "reseed.db")
	if err := InitDB(databasePath, []byte(testSeed)); err != nil {
		t.Fatal(err)
	}
	const deviceID = "b5bc0a48-04fc-4e0c-bf81-b1fb0cbf63d7"
	if _, _, err := ToggleVote("one", deviceID); err != nil {
		t.Fatal(err)
	}
	if err := Close(); err != nil {
		t.Fatal(err)
	}
	if err := InitDB(databasePath, []byte(testSeed)); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = Close() })

	items, err := GetAllVisualizers(deviceID)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].VoteCount != 1 || !items[0].UserVoted {
		t.Fatalf("vote was not preserved after reseed: %#v", items)
	}
}

func TestLegacyPopularityColumnIsRemoved(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "migration.db")
	if err := InitDB(databasePath, []byte(testSeed)); err != nil {
		t.Fatal(err)
	}
	if _, err := DB.Exec("ALTER TABLE visualizers ADD COLUMN popularity_rank INTEGER NOT NULL DEFAULT 0"); err != nil {
		t.Fatal(err)
	}
	if err := Close(); err != nil {
		t.Fatal(err)
	}

	if err := InitDB(databasePath, []byte(testSeed)); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = Close() })

	exists, err := columnExists("visualizers", "popularity_rank")
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("legacy popularity_rank column still exists")
	}
}
