package main

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStaticHandlerServesRoutesAndRejectsUnknownPaths(t *testing.T) {
	staticFS, err := fs.Sub(wwwEmbedFS, "www")
	if err != nil {
		t.Fatal(err)
	}
	handler := staticHandler(staticFS)

	for _, test := range []struct {
		name           string
		path           string
		status         int
		contentType    string
		acceptEncoding string
	}{
		{name: "home", path: "/", status: http.StatusOK, contentType: "text/html"},
		{
			name:           "brotli home",
			path:           "/",
			status:         http.StatusOK,
			contentType:    "text/html",
			acceptEncoding: "br",
		},
		{name: "privacy", path: "/privacy", status: http.StatusOK, contentType: "text/html"},
		{
			name:        "web manifest",
			path:        "/site.webmanifest",
			status:      http.StatusOK,
			contentType: "application/manifest+json",
		},
		{name: "missing", path: "/missing-page", status: http.StatusNotFound, contentType: "text/html"},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			if test.acceptEncoding != "" {
				request.Header.Set("Accept-Encoding", test.acceptEncoding)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status: got %d, want %d", response.Code, test.status)
			}
			if contentType := response.Header().Get("Content-Type"); !strings.HasPrefix(contentType, test.contentType) {
				t.Fatalf("content type: got %q, want prefix %q", contentType, test.contentType)
			}
			if test.contentType == "text/html" &&
				!strings.Contains(response.Header().Get("Content-Security-Policy"), "'sha256-") {
				t.Fatal("HTML response CSP does not allow its inline hydration scripts by hash")
			}
			if test.acceptEncoding != "" &&
				response.Header().Get("Content-Encoding") != test.acceptEncoding {
				t.Fatalf(
					"content encoding: got %q, want %q",
					response.Header().Get("Content-Encoding"),
					test.acceptEncoding,
				)
			}
		})
	}
}

func TestSecurityHeaders(t *testing.T) {
	handler := securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	for _, name := range []string{
		"Content-Security-Policy",
		"Permissions-Policy",
		"Referrer-Policy",
		"X-Content-Type-Options",
		"X-Frame-Options",
	} {
		if response.Header().Get(name) == "" {
			t.Errorf("missing %s", name)
		}
	}
}
