package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSPreflight(t *testing.T) {
	for _, tc := range []struct {
		name, origin, method, header string
		allowed                      bool
	}{
		{"allowed update", "http://localhost:3000", http.MethodPut, "content-type", true},
		{"allowed delete", "http://localhost:3000", http.MethodDelete, "", true},
		{"unknown origin", "https://other.example", http.MethodPut, "content-type", false},
		{"unknown method", "http://localhost:3000", http.MethodPatch, "", false},
		{"unknown header", "http://localhost:3000", http.MethodPut, "x-custom", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			handler := WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusUnauthorized)
			}), []string{"http://localhost:3000"})
			r := httptest.NewRequest(http.MethodOptions, "/api/task", nil)
			r.Header.Set("Origin", tc.origin)
			r.Header.Set("Access-Control-Request-Method", tc.method)
			if tc.header != "" {
				r.Header.Set("Access-Control-Request-Headers", tc.header)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			if called || w.Code != http.StatusNoContent {
				t.Fatalf("preflight reached application: called=%v, status=%d", called, w.Code)
			}
			allowed := w.Header().Get("Access-Control-Allow-Origin") == tc.origin
			if allowed != tc.allowed {
				t.Fatalf("origin allowed = %v, want %v", allowed, tc.allowed)
			}
			if tc.allowed && (w.Header().Get("Access-Control-Allow-Methods") != tc.method || w.Header().Get("Access-Control-Allow-Credentials") != "true") {
				t.Fatalf("incomplete preflight headers: %v", w.Header())
			}
		})
	}
}

func TestCORSKeepsApplicationResponse(t *testing.T) {
	for _, origin := range []string{"", "http://localhost:3000", "https://other.example"} {
		t.Run(origin, func(t *testing.T) {
			called := false
			handler := WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusUnauthorized)
			}), []string{"http://localhost:3000"})
			r := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
			if origin != "" {
				r.Header.Set("Origin", origin)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			if !called || w.Code != http.StatusUnauthorized {
				t.Fatal("CORS middleware bypassed application response")
			}
			wantOrigin := ""
			if origin == "http://localhost:3000" {
				wantOrigin = origin
			}
			if w.Header().Get("Access-Control-Allow-Origin") != wantOrigin {
				t.Fatalf("unexpected allow origin: %v", w.Header())
			}
		})
	}
}
