package runner

import "testing"

func TestParseWtxAddOutput(t *testing.T) {
	raw := []byte(`{"name":"feature","branch":"feature","path":"/tmp/repo/worktrees/feature"}`)

	out, err := parseWtxAddOutput(raw)
	if err != nil {
		t.Fatalf("parseWtxAddOutput returned error: %v", err)
	}
	if out.Path != "/tmp/repo/worktrees/feature" {
		t.Fatalf("path mismatch: got %q", out.Path)
	}
}

func TestParseWtxAddOutputMissingPath(t *testing.T) {
	_, err := parseWtxAddOutput([]byte(`{"name":"feature"}`))
	if err == nil {
		t.Fatal("expected error when path is missing")
	}
}

func TestNewAllowsNonRepoPath(t *testing.T) {
	storeDir := t.TempDir()
	nonRepo := t.TempDir()

	r, err := New(nonRepo, storeDir)
	if err != nil {
		t.Fatalf("New returned error for non-repo path: %v", err)
	}
	if r == nil {
		t.Fatal("expected runner instance")
	}
}
