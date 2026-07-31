package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Shik3i/KoalaGithub/internal/db"
)

type GitHubRepoResponse struct {
	StargazersCount int `json:"stargazers_count"`
}

// ExtractOwnerAndRepo extracts (owner, repo) from a GitHub repository URL
func ExtractOwnerAndRepo(repoURL string) (string, string, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid github repo path: %s", u.Path)
	}
	return parts[0], parts[1], nil
}

// FetchGitHubStars fetches the current stargazers_count from GitHub REST API
func FetchGitHubStars(owner, repo string) (int, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "KoalaGitHub-StarsFetcher/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var data GitHubRepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.StargazersCount, nil
}

// StartStaggeredStarsRefresher runs a background goroutine that refreshes one visualizer at a time
// spaced out evenly over 24 hours to avoid GitHub API rate limits.
func StartStaggeredStarsRefresher() {
	go func() {
		// Wait 5 seconds after server startup before initial check
		time.Sleep(5 * time.Second)

		for {
			visualizers, err := db.GetAllVisualizers("")
			if err != nil || len(visualizers) == 0 {
				time.Sleep(5 * time.Minute)
				continue
			}

			now := time.Now()
			// Calculate interval to space requests evenly across 24 hours (86400 seconds)
			interval := time.Duration(86400/len(visualizers)) * time.Second
			if interval < 10*time.Second {
				interval = 10 * time.Second
			}

			for idx, item := range visualizers {
				owner, repo, err := ExtractOwnerAndRepo(item.RepositoryURL)
				if err != nil {
					continue
				}

				stars, err := FetchGitHubStars(owner, repo)
				if err != nil {
					log.Printf("[Stars Refresher] Could not fetch stars for %s/%s: %v", owner, repo, err)
				} else {
					lastRefreshed := now.Format(time.RFC3339)
					nextRefresh := now.Add(24 * time.Hour).Format(time.RFC3339)
					if err := db.UpdateVisualizerStars(item.ID, stars, lastRefreshed, nextRefresh); err != nil {
						log.Printf("[Stars Refresher] Error updating stars for %s: %v", item.ID, err)
					} else {
						log.Printf("[Stars Refresher] Updated %s (%s/%s) -> ⭐ %d stars", item.ID, owner, repo, stars)
					}
				}

				// Stagger: sleep interval between each repository fetch
				if idx < len(visualizers)-1 {
					time.Sleep(interval)
				}
			}

			// Sleep for remaining day before starting next round
			time.Sleep(1 * time.Hour)
		}
	}()
}
