package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Shik3i/KoalaGithub/internal/db"
	"github.com/Shik3i/KoalaGithub/internal/models"
)

var githubClient = &http.Client{Timeout: 10 * time.Second}

type GitHubRepoResponse struct {
	StargazersCount int `json:"stargazers_count"`
}

func ExtractOwnerAndRepo(repoURL string) (string, string, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}
	if !strings.EqualFold(u.Scheme, "https") || !strings.EqualFold(u.Hostname(), "github.com") {
		return "", "", fmt.Errorf("unsupported repository host: %s", u.Hostname())
	}
	parts := strings.Split(strings.Trim(strings.TrimSuffix(u.Path, ".git"), "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid GitHub repository path: %s", u.Path)
	}
	return parts[0], parts[1], nil
}

func FetchGitHubStars(ctx context.Context, owner, repo string) (int, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", url.PathEscape(owner), url.PathEscape(repo))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "KoalaGitHub/1.0")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := githubClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var data GitHubRepoResponse
	decoder := json.NewDecoder(io.LimitReader(resp.Body, 1<<20))
	if err := decoder.Decode(&data); err != nil {
		return 0, err
	}
	return data.StargazersCount, nil
}

func StartStaggeredStarsRefresher(ctx context.Context) {
	go func() {
		if !wait(ctx, 5*time.Second) {
			return
		}
		for {
			visualizers, err := db.GetAllVisualizers("")
			if err != nil || len(visualizers) == 0 {
				if !wait(ctx, 5*time.Minute) {
					return
				}
				continue
			}
			refreshUniqueRepositories(ctx, visualizers)
			if !wait(ctx, 24*time.Hour) {
				return
			}
		}
	}()
}

func refreshUniqueRepositories(ctx context.Context, visualizers []models.Visualizer) {
	groups := make(map[string][]models.Visualizer)
	order := make([]string, 0, len(visualizers))
	for _, item := range visualizers {
		if refreshedAt, err := time.Parse(time.RFC3339, item.LastRefreshedAt); err == nil &&
			time.Now().Before(refreshedAt.Add(23*time.Hour)) {
			continue
		}
		owner, repo, err := ExtractOwnerAndRepo(item.RepositoryURL)
		if err != nil {
			log.Printf("[Stars] Skip %s: %v", item.ID, err)
			continue
		}
		key := strings.ToLower(owner + "/" + repo)
		if _, exists := groups[key]; !exists {
			order = append(order, key)
		}
		groups[key] = append(groups[key], item)
	}
	if len(order) == 0 {
		return
	}
	for index, key := range order {
		parts := strings.SplitN(key, "/", 2)
		stars, err := FetchGitHubStars(ctx, parts[0], parts[1])
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("[Stars] Fetch %s: %v", key, err)
		} else {
			refreshedAt := time.Now().UTC()
			nextRefresh := refreshedAt.Add(24 * time.Hour)
			for _, item := range groups[key] {
				if err := db.UpdateVisualizerStars(item.ID, stars, refreshedAt.Format(time.RFC3339), nextRefresh.Format(time.RFC3339)); err != nil {
					log.Printf("[Stars] Update %s: %v", item.ID, err)
				}
			}
		}
		if index < len(order)-1 && !wait(ctx, 2*time.Second) {
			return
		}
	}
}

func wait(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
