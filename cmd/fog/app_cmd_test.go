package main

import (
	"runtime"
	"strings"
	"testing"
)

func TestRunAppUnsupportedPlatformGuard(t *testing.T) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		t.Skip("platform guard only relevant on unsupported OS")
	}
	err := runApp()
	if err == nil {
		t.Fatal("expected unsupported platform error")
	}
	if !strings.Contains(err.Error(), "supported on macOS and Linux") {
		t.Fatalf("unexpected error: %v", err)
	}
}
