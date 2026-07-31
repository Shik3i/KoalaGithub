package main

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/Shik3i/KoalaGithub/internal/db"
	"github.com/Shik3i/KoalaGithub/internal/handlers"
	"github.com/Shik3i/KoalaGithub/internal/services"
)

var version = "dev"
var inlineScriptPattern = regexp.MustCompile(`(?s)<script(?:\s[^>]*)?>(.*?)</script>`)

//go:embed all:www
var wwwEmbedFS embed.FS

//go:embed src/lib/data/visualizers.json
var seedJSON []byte

func main() {
	port := envOrDefault("PORT", "8080")
	dbPath := envOrDefault("DB_PATH", "./data/koalagithub.db")
	if err := db.InitDB(dbPath, seedJSON); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	services.StartStaggeredStarsRefresher(ctx)
	handlers.AppVersion = version

	staticFS, err := fs.Sub(wwwEmbedFS, "www")
	if err != nil {
		log.Fatalf("Embedded frontend unavailable: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handlers.HealthHandler)
	mux.HandleFunc("GET /api/visualizers", handlers.GetVisualizersHandler)
	mux.HandleFunc("POST /api/visualizers/{id}/vote", handlers.VoteHandler)
	mux.HandleFunc("DELETE /api/votes", handlers.DeleteVotesHandler)
	mux.Handle("/", staticHandler(staticFS))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           securityHeaders(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown: %v", err)
		}
	}()

	log.Printf("KoalaGitHub %s starting on http://localhost:%s", version, port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func staticHandler(staticFS fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		cleanPath := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
		if cleanPath == "." || cleanPath == "" {
			cleanPath = "index.html"
		}
		candidates := []string{cleanPath}
		if path.Ext(cleanPath) == "" {
			candidates = append(candidates, cleanPath+".html", path.Join(cleanPath, "index.html"))
		}
		for _, candidate := range candidates {
			if serveStaticFile(w, r, staticFS, candidate, http.StatusOK) {
				return
			}
		}
		if serveStaticFile(w, r, staticFS, "404.html", http.StatusNotFound) {
			return
		}
		http.NotFound(w, r)
	})
}

func serveStaticFile(w http.ResponseWriter, r *http.Request, staticFS fs.FS, name string, status int) bool {
	encoding := ""
	servedName := name
	acceptEncoding := r.Header.Get("Accept-Encoding")
	if strings.Contains(acceptEncoding, "br") && fileExists(staticFS, name+".br") {
		encoding, servedName = "br", name+".br"
	} else if strings.Contains(acceptEncoding, "gzip") && fileExists(staticFS, name+".gz") {
		encoding, servedName = "gzip", name+".gz"
	} else if !fileExists(staticFS, name) {
		return false
	}

	content, err := fs.ReadFile(staticFS, servedName)
	if err != nil {
		return false
	}
	if encoding != "" {
		w.Header().Set("Content-Encoding", encoding)
		w.Header().Add("Vary", "Accept-Encoding")
	}
	contentType := mime.TypeByExtension(path.Ext(name))
	if path.Ext(name) == ".webmanifest" {
		contentType = "application/manifest+json"
	}
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	if path.Ext(name) == ".html" {
		html := content
		if encoding != "" {
			uncompressed, err := fs.ReadFile(staticFS, name)
			if err != nil {
				return false
			}
			html = uncompressed
		}
		w.Header().Set("Content-Security-Policy", contentSecurityPolicy(html))
	}
	if strings.HasPrefix(name, "_app/immutable/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}
	w.WriteHeader(status)
	if r.Method != http.MethodHead {
		_, _ = w.Write(content)
	}
	return true
}

func fileExists(staticFS fs.FS, name string) bool {
	info, err := fs.Stat(staticFS, name)
	return err == nil && !info.IsDir()
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", contentSecurityPolicy(nil))
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

func contentSecurityPolicy(html []byte) string {
	scriptSources := []string{"'self'"}
	for _, match := range inlineScriptPattern.FindAllSubmatch(html, -1) {
		digest := sha256.Sum256(match[1])
		scriptSources = append(scriptSources, "'sha256-"+base64.StdEncoding.EncodeToString(digest[:])+"'")
	}
	return "default-src 'self'; img-src 'self' data: https:; style-src 'self' 'unsafe-inline'; script-src " +
		strings.Join(scriptSources, " ") +
		"; connect-src 'self'; font-src 'self'; object-src 'none'; base-uri 'self'; frame-ancestors 'none'; form-action 'self'; upgrade-insecure-requests"
}
