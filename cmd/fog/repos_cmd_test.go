package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/darkLord19/foglet/internal/ghcli"
)

func TestParseIndex(t *testing.T) {
	idx, err := parseIndex("2", 3)
	if err != nil {
		t.Fatalf("parseIndex returned error: %v", err)
	}
	if idx != 1 {
		t.Fatalf("parseIndex mismatch: got %d want 1", idx)
	}
}

func TestParseIndexesDedup(t *testing.T) {
	indexes, err := parseIndexes("1, 2, 2, 3", 3)
	if err != nil {
		t.Fatalf("parseIndexes returned error: %v", err)
	}
	want := []int{0, 1, 2}
	if !reflect.DeepEqual(indexes, want) {
		t.Fatalf("parseIndexes mismatch: got %v want %v", indexes, want)
	}
}

func TestSplitRepoFullName(t *testing.T) {
	owner, name, err := splitRepoFullName("acme/api")
	if err != nil {
		t.Fatalf("splitRepoFullName returned error: %v", err)
	}
	if owner != "acme" || name != "api" {
		t.Fatalf("unexpected split result: %q/%q", owner, name)
	}

	if _, _, err := splitRepoFullName("../repo"); err == nil {
		t.Fatalf("expected invalid segment error")
	}
}

func TestSelectReposByFullName(t *testing.T) {
	repos := []ghcli.Repo{
		{NameWithOwner: "acme/api", Name: "api"},
		{NameWithOwner: "acme/web", Name: "web"},
	}
	selected, err := selectRepos(repos, "acme/web")
	if err != nil {
		t.Fatalf("selectRepos returned error: %v", err)
	}
	if len(selected) != 1 || selected[0].NameWithOwner != "acme/web" {
		t.Fatalf("unexpected selected repos: %+v", selected)
	}
}

func TestEnsureBareRepoInitialized(t *testing.T) {
	tmp := t.TempDir()
	barePath := filepath.Join(tmp, "repo.git")
	basePath := filepath.Join(tmp, "base")

	repo := ghcli.Repo{
		NameWithOwner: "acme/api",
	}

	var calls []string
	origClone := cloneGhRepoFn
	origRunner := gitRunner
	t.Cleanup(func() {
		cloneGhRepoFn = origClone
		gitRunner = origRunner
	})

	cloneGhRepoFn = func(fullName, destPath string) error {
		calls = append(calls, "clone "+fullName)
		return os.MkdirAll(destPath, 0o755)
	}

	gitRunner = func(extraEnv []string, args ...string) error {
		_ = extraEnv
		calls = append(calls, strings.Join(args, " "))
		if len(args) >= 5 && args[0] == "--git-dir" {
			return os.MkdirAll(basePath, 0o755)
		}
		return nil
	}

	if err := ensureBareRepoInitialized(repo, barePath, basePath); err != nil {
		t.Fatalf("ensureBareRepoInitialized failed: %v", err)
	}

	if len(calls) != 2 {
		t.Fatalf("expected 2 calls (clone + worktree add), got %d: %v", len(calls), calls)
	}
	if calls[0] != "clone acme/api" {
		t.Fatalf("unexpected clone call: %q", calls[0])
	}
	if !strings.Contains(calls[1], "worktree add") {
		t.Fatalf("expected worktree add call, got %q", calls[1])
	}
}
