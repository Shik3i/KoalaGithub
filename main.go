package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Shik3i/KoalaGithub/internal/db"
	"github.com/Shik3i/KoalaGithub/internal/handlers"
)

//go:embed all:www
var wwwEmbedFS embed.FS

//go:embed src/lib/data/visualizers.json
var seedJSON []byte

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/koalagithub.db"
	}

	if err := db.InitDB(dbPath, seedJSON); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	mux := http.NewServeMux()

	// API Endpoints
	mux.HandleFunc("GET /api/health", handlers.HealthHandler)
	mux.HandleFunc("GET /api/visualizers", handlers.GetVisualizersHandler)
	mux.HandleFunc("POST /api/visualizers/{id}/vote", handlers.VoteHandler)

	// Static Web Assets & SPA Routing
	staticFS, err := fs.Sub(wwwEmbedFS, "www")
	if err != nil {
		log.Printf("Notice: Embedded www folder not found, attempting local disk fallback")
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Do not match API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		cleanPath := filepath.Clean(r.URL.Path)

		// Serve file if it exists in staticFS or disk
		if staticFS != nil {
			f, err := staticFS.Open(strings.TrimPrefix(cleanPath, "/"))
			if err == nil {
				f.Close()
				http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
				return
			}
		} else {
			diskPath := filepath.Join("www", cleanPath)
			if info, err := os.Stat(diskPath); err == nil && !info.IsDir() {
				http.ServeFile(w, r, diskPath)
				return
			}
		}

		// Fallback to index.html for SPA/SSG routing
		if staticFS != nil {
			r.URL.Path = "/"
			http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
		} else {
			http.ServeFile(w, r, filepath.Join("www", "index.html"))
		}
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("KoalaGitHub server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
