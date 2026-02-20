package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestSocketModeRunOnceAppMentionParseErrorPostsThreadMessage(t *testing.T) {
	t.Parallel()

	upgrader := websocket.Upgrader{}
	ackCh := make(chan map[string]string, 1)
	chatCh := make(chan map[string]string, 1)

	var wsURL string
	mux := http.NewServeMux()

	mux.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer xapp-test" {
			t.Errorf("unexpected app token header: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":  true,
			"url": wsURL,
		})
	})

	mux.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		envelope := map[string]any{
			"envelope_id": "env-1",
			"type":        "events_api",
			"payload": map[string]any{
				"event": map[string]any{
					"type":    "app_mention",
					"text":    "<@U123> missing options block",
					"channel": "C123",
					"ts":      "111.222",
				},
			},
		}
		if err := conn.WriteJSON(envelope); err != nil {
			t.Errorf("write envelope failed: %v", err)
			return
		}

		_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		var ack map[string]string
		if err := conn.ReadJSON(&ack); err != nil {
			t.Errorf("read ack failed: %v", err)
			return
		}
		ackCh <- ack
	})

	mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer xoxb-test" {
			t.Errorf("unexpected bot token header: %q", got)
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode postMessage payload failed: %v", err)
			http.Error(w, "bad payload", http.StatusBadRequest)
			return
		}
		chatCh <- payload

		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok": true,
			"ts": "111.223",
		})
	})

	server := newHTTPTestServerOrSkip(t, mux)
	defer server.Close()
	wsURL = toWSURL(server.URL) + "/socket"

	sm := NewSocketMode(nil, nil, "xapp-test", "xoxb-test")
	sm.httpClient = server.Client()
	sm.connectionsOpenURL = server.URL + "/open"
	sm.postMessageURL = server.URL + "/chat.postMessage"
	sm.dialer = websocket.DefaultDialer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = sm.runOnce(ctx)

	select {
	case ack := <-ackCh:
		if ack["envelope_id"] != "env-1" {
			t.Fatalf("unexpected ack envelope_id: %q", ack["envelope_id"])
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for websocket ack")
	}

	select {
	case payload := <-chatCh:
		if payload["channel"] != "C123" {
			t.Fatalf("unexpected channel: %q", payload["channel"])
		}
		if payload["thread_ts"] != "111.222" {
			t.Fatalf("unexpected thread_ts: %q", payload["thread_ts"])
		}
		if !strings.Contains(payload["text"], "options block is required") {
			t.Fatalf("unexpected error text: %q", payload["text"])
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for chat.postMessage payload")
	}
}

func TestSocketModeRunOnceSlashCommandParseErrorPostsResponseURL(t *testing.T) {
	t.Parallel()

	upgrader := websocket.Upgrader{}
	ackCh := make(chan map[string]string, 1)
	responseCh := make(chan map[string]any, 1)

	var wsURL string
	mux := http.NewServeMux()

	mux.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":  true,
			"url": wsURL,
		})
	})

	mux.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		envelope := map[string]any{
			"envelope_id": "env-2",
			"type":        "slash_commands",
			"payload": map[string]any{
				"channel_id":   "C123",
				"text":         "invalid text",
				"response_url": serverURLForTest(r) + "/response",
			},
		}
		if err := conn.WriteJSON(envelope); err != nil {
			t.Errorf("write envelope failed: %v", err)
			return
		}

		_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		var ack map[string]string
		if err := conn.ReadJSON(&ack); err != nil {
			t.Errorf("read ack failed: %v", err)
			return
		}
		ackCh <- ack
	})

	mux.HandleFunc("/response", func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode response_url payload failed: %v", err)
			http.Error(w, "bad payload", http.StatusBadRequest)
			return
		}
		responseCh <- payload
		_, _ = w.Write([]byte(`ok`))
	})

	server := newHTTPTestServerOrSkip(t, mux)
	defer server.Close()
	wsURL = toWSURL(server.URL) + "/socket"

	sm := NewSocketMode(nil, nil, "xapp-test", "xoxb-test")
	sm.httpClient = server.Client()
	sm.connectionsOpenURL = server.URL + "/open"
	sm.postMessageURL = server.URL + "/chat.postMessage"
	sm.dialer = websocket.DefaultDialer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = sm.runOnce(ctx)

	select {
	case ack := <-ackCh:
		if ack["envelope_id"] != "env-2" {
			t.Fatalf("unexpected ack envelope_id: %q", ack["envelope_id"])
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for websocket ack")
	}

	select {
	case payload := <-responseCh:
		text, _ := payload["text"].(string)
		responseType, _ := payload["response_type"].(string)
		if responseType != "ephemeral" {
			t.Fatalf("unexpected response_type: %q", responseType)
		}
		if !strings.Contains(text, "options block is required") {
			t.Fatalf("unexpected response text: %q", text)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for response_url payload")
	}
}

func toWSURL(httpURL string) string {
	return "ws" + strings.TrimPrefix(httpURL, "http")
}

func serverURLForTest(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

func newHTTPTestServerOrSkip(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("skipping socket integration test: %v", r)
		}
	}()
	return httptest.NewServer(handler)
}
