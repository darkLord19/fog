package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darkLord19/wtx/internal/runner"
	"github.com/darkLord19/wtx/internal/state"
)

func TestHandleCreateTaskRequiresRepo(t *testing.T) {
	srv := newTestServer(t)

	body := []byte(`{"branch":"feature-a","prompt":"do thing"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	srv.handleCreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleCreateTaskRejectsUnknownRepo(t *testing.T) {
	srv := newTestServer(t)

	body := []byte(`{"repo":"missing","branch":"feature-a","prompt":"do thing"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	srv.handleCreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d want %d", w.Code, http.StatusBadRequest)
	}
}

func newTestServer(t *testing.T) *Server {
	t.Helper()

	cfgDir := t.TempDir()
	r, err := runner.New("", cfgDir)
	if err != nil {
		t.Fatalf("new runner failed: %v", err)
	}

	st, err := state.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("new state store failed: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	if err := st.SetDefaultTool("claude"); err != nil {
		t.Fatalf("set default tool failed: %v", err)
	}

	return New(r, st, 8080)
}
