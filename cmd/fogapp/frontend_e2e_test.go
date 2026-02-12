package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

func TestDesktopFrontendSmokeFlows(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e smoke in short mode")
	}

	chromePath := findChromeBinary()
	if chromePath == "" {
		t.Skip("skipping desktop e2e smoke: no Chrome/Chromium binary found")
	}

	mockAPI := newMockFogAPI()
	apiServer := httptest.NewServer(mockAPI)
	defer apiServer.Close()

	frontendServer := newFrontendHarnessServer(t, apiServer.URL)
	defer frontendServer.Close()

	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("window-size", "1400,1000"),
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 55*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(frontendServer.URL),
		chromedp.WaitVisible("#new-session-form", chromedp.ByQuery),
		waitTextContains("#daemon-badge", "connected"),
		waitTextContains("#timeline-summary", "owner/repo"),
		waitTextContains("#timeline-runs", "COMPLETED"),
		waitTextContains("#timeline-events", "Run completed"),
		waitTextContains("#timeline-actions", "Open Branch"),

		chromedp.SetValue("#new-prompt", "Implement desktop smoke flow", chromedp.ByQuery),
		chromedp.Click("#new-submit", chromedp.ByQuery),
		waitTextContains("#new-status", "Queued session"),

		chromedp.SetValue("#followup-prompt", "Add regression tests", chromedp.ByQuery),
		chromedp.Click("#followup-submit", chromedp.ByQuery),
		waitTextContains("#followup-status", "Queued run"),

		chromedp.Click("#discover-btn", chromedp.ByQuery),
		chromedp.WaitVisible("#repo-0", chromedp.ByQuery),
		chromedp.Click("#repo-0", chromedp.ByQuery),
		chromedp.Click("#import-btn", chromedp.ByQuery),
		waitTextContains("#repos-status", "Imported"),

		chromedp.SetValue("#settings-prefix", "team", chromedp.ByQuery),
		chromedp.SetValue("#settings-pat", "ghp_mock_token", chromedp.ByQuery),
		chromedp.Click("#settings-submit", chromedp.ByQuery),
		waitTextContains("#settings-status", "Saved"),

		chromedp.Click("#cloud-save", chromedp.ByQuery),
		waitTextContains("#cloud-status", "Cloud URL saved"),
	)
	if err != nil {
		t.Fatalf("desktop e2e flow failed: %v", err)
	}

	stats := mockAPI.stats()
	if stats.createSessionCount < 1 {
		t.Fatalf("expected create session call, got %+v", stats)
	}
	if stats.followupCount < 1 {
		t.Fatalf("expected follow-up call, got %+v", stats)
	}
	if stats.discoverCount < 1 || stats.importCount < 1 {
		t.Fatalf("expected repo discover/import calls, got %+v", stats)
	}
	if stats.settingsPutCount < 1 {
		t.Fatalf("expected settings update call, got %+v", stats)
	}
	if stats.cloudPutCount < 1 {
		t.Fatalf("expected cloud url save call, got %+v", stats)
	}
}

type e2eStats struct {
	createSessionCount int
	followupCount      int
	discoverCount      int
	importCount        int
	settingsPutCount   int
	cloudPutCount      int
}

type mockFogAPI struct {
	mu sync.Mutex

	counters e2eStats

	settings map[string]interface{}
	repos    []map[string]interface{}
	cloud    map[string]interface{}
	sessions []map[string]interface{}
	runs     map[string][]map[string]interface{}
	events   map[string][]map[string]interface{}
}

func newMockFogAPI() *mockFogAPI {
	now := time.Now().UTC()
	sessions := []map[string]interface{}{
		{
			"id":         "session-1",
			"repo_name":  "owner/repo",
			"branch":     "fog/session-one",
			"tool":       "claude",
			"status":     "COMPLETED",
			"busy":       false,
			"autopr":     false,
			"pr_url":     "",
			"updated_at": now.Format(time.RFC3339),
		},
	}
	runs := map[string][]map[string]interface{}{
		"session-1": {
			{
				"id":         "run-1",
				"session_id": "session-1",
				"prompt":     "Initial prompt",
				"state":      "COMPLETED",
				"created_at": now.Add(-3 * time.Minute).Format(time.RFC3339),
				"updated_at": now.Add(-2 * time.Minute).Format(time.RFC3339),
			},
		},
	}
	events := map[string][]map[string]interface{}{
		"run-1": {
			{"id": 1, "run_id": "run-1", "ts": now.Add(-3 * time.Minute).Format(time.RFC3339), "type": "ai_start", "message": "Running AI tool"},
			{"id": 2, "run_id": "run-1", "ts": now.Add(-2 * time.Minute).Format(time.RFC3339), "type": "complete", "message": "Run completed"},
		},
	}

	return &mockFogAPI{
		settings: map[string]interface{}{
			"default_tool":        "claude",
			"branch_prefix":       "fog",
			"has_github_token":    true,
			"onboarding_required": false,
			"available_tools":     []string{"claude", "cursor"},
		},
		repos: []map[string]interface{}{
			{
				"name":               "owner/repo",
				"url":                "https://github.com/owner/repo.git",
				"default_branch":     "main",
				"base_worktree_path": "/tmp/owner-repo/base",
			},
		},
		cloud: map[string]interface{}{
			"cloud_url":        "https://fog-cloud.example",
			"device_id":        "device-1",
			"has_device_token": true,
			"paired":           true,
		},
		sessions: sessions,
		runs:     runs,
		events:   events,
	}
}

func (m *mockFogAPI) statsSnapshot() e2eStats {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.counters
}

func (m *mockFogAPI) stats() e2eStats {
	return m.statsSnapshot()
}

func (m *mockFogAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	writeJSON := func(code int, payload interface{}) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(payload)
	}

	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/api/settings":
		writeJSON(http.StatusOK, m.settings)
		return
	case r.Method == http.MethodPut && r.URL.Path == "/api/settings":
		m.counters.settingsPutCount++
		var in map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&in)
		if v, ok := in["default_tool"].(string); ok && strings.TrimSpace(v) != "" {
			m.settings["default_tool"] = v
		}
		if v, ok := in["branch_prefix"].(string); ok && strings.TrimSpace(v) != "" {
			m.settings["branch_prefix"] = v
		}
		if v, ok := in["github_pat"].(string); ok && strings.TrimSpace(v) != "" {
			m.settings["has_github_token"] = true
		}
		writeJSON(http.StatusOK, m.settings)
		return
	case r.Method == http.MethodGet && r.URL.Path == "/api/repos":
		writeJSON(http.StatusOK, m.repos)
		return
	case r.Method == http.MethodPost && r.URL.Path == "/api/repos/discover":
		m.counters.discoverCount++
		writeJSON(http.StatusOK, []map[string]interface{}{
			{"full_name": "owner/new-repo", "default_branch": "main"},
		})
		return
	case r.Method == http.MethodPost && r.URL.Path == "/api/repos/import":
		m.counters.importCount++
		writeJSON(http.StatusOK, map[string]interface{}{
			"imported": []string{"owner/new-repo"},
			"failed":   []string{},
		})
		return
	case r.Method == http.MethodGet && r.URL.Path == "/api/cloud":
		writeJSON(http.StatusOK, m.cloud)
		return
	case r.Method == http.MethodPut && r.URL.Path == "/api/cloud":
		m.counters.cloudPutCount++
		var in map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&in)
		if v, ok := in["cloud_url"].(string); ok && strings.TrimSpace(v) != "" {
			m.cloud["cloud_url"] = v
		}
		writeJSON(http.StatusOK, m.cloud)
		return
	case r.Method == http.MethodPost && r.URL.Path == "/api/cloud/pair":
		m.cloud["paired"] = true
		writeJSON(http.StatusOK, m.cloud)
		return
	case r.Method == http.MethodPost && r.URL.Path == "/api/cloud/unpair":
		m.cloud["paired"] = false
		writeJSON(http.StatusOK, m.cloud)
		return
	case r.Method == http.MethodGet && r.URL.Path == "/api/sessions":
		writeJSON(http.StatusOK, m.sessions)
		return
	case r.Method == http.MethodPost && r.URL.Path == "/api/sessions":
		m.counters.createSessionCount++
		id := "session-" + strconv.Itoa(len(m.sessions)+1)
		runID := "run-" + strconv.Itoa(len(m.events)+1)
		now := time.Now().UTC().Format(time.RFC3339)
		session := map[string]interface{}{
			"id":         id,
			"repo_name":  "owner/repo",
			"branch":     "fog/new-session",
			"tool":       "claude",
			"status":     "CREATED",
			"busy":       true,
			"autopr":     false,
			"pr_url":     "",
			"updated_at": now,
		}
		m.sessions = append([]map[string]interface{}{session}, m.sessions...)
		m.runs[id] = []map[string]interface{}{
			{
				"id":         runID,
				"session_id": id,
				"prompt":     "New session",
				"state":      "CREATED",
				"created_at": now,
				"updated_at": now,
			},
		}
		m.events[runID] = []map[string]interface{}{
			{"id": 1, "run_id": runID, "ts": now, "type": "setup", "message": "queued"},
		}
		writeJSON(http.StatusAccepted, map[string]interface{}{
			"session_id": id,
			"run_id":     runID,
			"status":     "accepted",
		})
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api/sessions/") {
		path := strings.TrimPrefix(r.URL.Path, "/api/sessions/")
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) == 1 && r.Method == http.MethodGet {
			sid := parts[0]
			session := map[string]interface{}{}
			for _, s := range m.sessions {
				if s["id"] == sid {
					session = s
					break
				}
			}
			writeJSON(http.StatusOK, map[string]interface{}{
				"session": session,
				"runs":    m.runs[sid],
			})
			return
		}
		if len(parts) == 2 && parts[1] == "runs" && r.Method == http.MethodPost {
			m.counters.followupCount++
			sid := parts[0]
			runID := "run-" + strconv.Itoa(len(m.events)+1)
			now := time.Now().UTC().Format(time.RFC3339)
			m.runs[sid] = append([]map[string]interface{}{
				{
					"id":         runID,
					"session_id": sid,
					"prompt":     "Follow-up",
					"state":      "CREATED",
					"created_at": now,
					"updated_at": now,
				},
			}, m.runs[sid]...)
			m.events[runID] = []map[string]interface{}{
				{"id": 1, "run_id": runID, "ts": now, "type": "ai_start", "message": "queued"},
			}
			writeJSON(http.StatusAccepted, map[string]interface{}{
				"run_id":  runID,
				"status":  "accepted",
				"session": sid,
			})
			return
		}
		if len(parts) == 4 && parts[1] == "runs" && parts[3] == "events" && r.Method == http.MethodGet {
			runID := parts[2]
			writeJSON(http.StatusOK, m.events[runID])
			return
		}
	}

	http.NotFound(w, r)
}

func newFrontendHarnessServer(t *testing.T, apiBaseURL string) *httptest.Server {
	t.Helper()

	root := filepath.Join("frontend")
	indexRaw, err := os.ReadFile(filepath.Join(root, "index.html"))
	if err != nil {
		t.Fatalf("read frontend index failed: %v", err)
	}
	appRaw, err := os.ReadFile(filepath.Join(root, "app.js"))
	if err != nil {
		t.Fatalf("read frontend app.js failed: %v", err)
	}
	cssRaw, err := os.ReadFile(filepath.Join(root, "styles.css"))
	if err != nil {
		t.Fatalf("read frontend styles failed: %v", err)
	}

	inject := "<script>window.__FOG_API_BASE_URL__ = " + strconv.Quote(apiBaseURL) + ";</script>\n  <script src=\"app.js\"></script>"
	index := strings.Replace(string(indexRaw), "<script src=\"app.js\"></script>", inject, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(index))
	})
	mux.HandleFunc("/app.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		_, _ = w.Write(appRaw)
	})
	mux.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(cssRaw)
	})

	return httptest.NewServer(mux)
}

func waitTextContains(selector, substring string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		deadline := time.Now().Add(12 * time.Second)
		for {
			var out string
			err := chromedp.Run(ctx, chromedp.Text(selector, &out, chromedp.ByQuery))
			if err == nil && strings.Contains(out, substring) {
				return nil
			}
			if time.Now().After(deadline) {
				if err != nil {
					return fmt.Errorf("waitTextContains(%s): %w", selector, err)
				}
				return fmt.Errorf("waitTextContains(%s): %q does not include %q", selector, out, substring)
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(120 * time.Millisecond):
			}
		}
	})
}

func findChromeBinary() string {
	candidates := []string{
		"google-chrome",
		"google-chrome-stable",
		"chromium",
		"chromium-browser",
	}
	for _, candidate := range candidates {
		if path, err := exec.LookPath(candidate); err == nil {
			return path
		}
	}
	return ""
}
