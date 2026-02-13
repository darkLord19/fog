package ai

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCommandPathFindsCommandInPATH(t *testing.T) {
	tempDir := t.TempDir()
	cmdName := "fog-test-cmd-path"
	cmdPath := filepath.Join(tempDir, cmdName)
	if err := os.WriteFile(cmdPath, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write command fixture failed: %v", err)
	}
	t.Setenv("PATH", tempDir)

	got := commandPath(cmdName)
	if got == "" {
		t.Fatalf("expected command path for %q", cmdName)
	}
	if got != cmdPath {
		t.Fatalf("unexpected command path: got %q want %q", got, cmdPath)
	}
}

func TestCommandPathFindsFallbackHomeLocalBin(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("fallback path test is unix-like specific")
	}

	homeDir := t.TempDir()
	fallbackDir := filepath.Join(homeDir, ".local", "bin")
	if err := os.MkdirAll(fallbackDir, 0o755); err != nil {
		t.Fatalf("mkdir fallback dir failed: %v", err)
	}

	cmdName := "fog-test-cmd-fallback"
	cmdPath := filepath.Join(fallbackDir, cmdName)
	if err := os.WriteFile(cmdPath, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write fallback command fixture failed: %v", err)
	}

	t.Setenv("HOME", homeDir)
	t.Setenv("PATH", "/usr/bin:/bin")

	got := commandPath(cmdName)
	if got != cmdPath {
		t.Fatalf("unexpected fallback command path: got %q want %q", got, cmdPath)
	}
}
