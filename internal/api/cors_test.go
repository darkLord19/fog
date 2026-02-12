package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCORSSetsHeaders(t *testing.T) {
	handler := WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("unexpected allow origin header: %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Fatal("expected allow methods header")
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected response code: %d", rec.Code)
	}
}

func TestWithCORSHandlesOptions(t *testing.T) {
	handler := WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run for OPTIONS")
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/tasks", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected response code: %d", rec.Code)
	}
}
