package daemon

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestIsHealthyWithClientOK(t *testing.T) {
	client := &http.Client{Transport: roundTrip(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("ok")),
			Header:     make(http.Header),
		}
	})}

	if !isHealthyWithClient(client, "http://example.local/health") {
		t.Fatal("expected health check to pass")
	}
}

func TestIsHealthyWithClientErrorStatus(t *testing.T) {
	client := &http.Client{Transport: roundTrip(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBufferString("err")),
			Header:     make(http.Header),
		}
	})}

	if isHealthyWithClient(client, "http://example.local/health") {
		t.Fatal("expected health check to fail")
	}
}

func TestWaitForHealthTimeout(t *testing.T) {
	start := time.Now()
	err := waitForHealth("http://127.0.0.1:1/health", 900*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if time.Since(start) < 900*time.Millisecond {
		t.Fatal("expected waitForHealth to wait for timeout duration")
	}
}

type roundTrip func(*http.Request) *http.Response

func (f roundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
