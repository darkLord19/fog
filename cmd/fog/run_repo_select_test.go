package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/darkLord19/foglet/internal/ghcli"
	"github.com/darkLord19/foglet/internal/state"
)

func TestResolveRepoNameForRunUsesFlag(t *testing.T) {
	origTTY := stdinIsTTYFn
	t.Cleanup(func() { stdinIsTTYFn = origTTY })
	stdinIsTTYFn = func() bool {
		t.Fatalf("stdinIsTTYFn should not be called when --repo is provided")
		return false
	}

	got, err := resolveRepoNameForRun(" acme/api ", nil)
	if err != nil {
		t.Fatalf("resolveRepoNameForRun returned error: %v", err)
	}
	if got != "acme/api" {
		t.Fatalf("resolveRepoNameForRun mismatch: got %q want %q", got, "acme/api")
	}
}

func TestResolveRepoNameForRunRejectsInvalidFlag(t *testing.T) {
	_, err := resolveRepoNameForRun("../repo", nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestResolveRepoNameForRunRequiresRepoWhenNotTTY(t *testing.T) {
	origTTY := stdinIsTTYFn
	t.Cleanup(func() { stdinIsTTYFn = origTTY })
	stdinIsTTYFn = func() bool { return false }

	_, err := resolveRepoNameForRun("", nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "--repo is required") {
		t.Fatalf("expected --repo required error, got %v", err)
	}
}

func TestResolveRepoNameForRunPromptsAndSelectsRepo(t *testing.T) {
	origTTY := stdinIsTTYFn
	origList := listGitHubReposFn
	origRead := readLineFn
	origAvail := isGhAvailableFn
	origAuth := isGhAuthenticatedFn
	t.Cleanup(func() {
		stdinIsTTYFn = origTTY
		listGitHubReposFn = origList
		readLineFn = origRead
		isGhAvailableFn = origAvail
		isGhAuthenticatedFn = origAuth
	})

	stdinIsTTYFn = func() bool { return true }
	isGhAvailableFn = func() bool { return true }
	isGhAuthenticatedFn = func() bool { return true }

	listGitHubReposFn = func() ([]ghcli.Repo, error) {
		return []ghcli.Repo{
			{NameWithOwner: "acme/api", Name: "api"},
			{NameWithOwner: "acme/web", Name: "web"},
		}, nil
	}
	readLineFn = func(prompt string) (string, error) {
		if !strings.Contains(prompt, "Select repository") {
			t.Fatalf("unexpected prompt: %q", prompt)
		}
		return "2", nil
	}

	got, err := resolveRepoNameForRun("", nil)
	if err != nil {
		t.Fatalf("resolveRepoNameForRun returned error: %v", err)
	}
	if got != "acme/web" {
		t.Fatalf("resolveRepoNameForRun mismatch: got %q want %q", got, "acme/web")
	}
}

func TestEnsureRepoRegisteredForRunReturnsExistingRepo(t *testing.T) {
	fogHome := t.TempDir()
	store, err := state.NewStore(fogHome)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	want := state.Repo{
		Name:             "acme/api",
		URL:              "https://github.com/acme/api.git",
		Host:             "github.com",
		Owner:            "acme",
		Repo:             "api",
		BarePath:         filepath.Join(fogHome, "bare.git"),
		BaseWorktreePath: filepath.Join(fogHome, "base"),
		DefaultBranch:    "main",
	}
	if _, err := store.UpsertRepo(want); err != nil {
		t.Fatalf("UpsertRepo failed: %v", err)
	}

	origList := listGitHubReposFn
	t.Cleanup(func() { listGitHubReposFn = origList })
	listGitHubReposFn = func() ([]ghcli.Repo, error) {
		t.Fatalf("listGitHubReposFn should not be called when repo already exists")
		return nil, nil
	}

	got, err := ensureRepoRegisteredForRun("acme/api", store, fogHome)
	if err != nil {
		t.Fatalf("ensureRepoRegisteredForRun returned error: %v", err)
	}
	if got.Name != want.Name {
		t.Fatalf("repo name mismatch: got %q want %q", got.Name, want.Name)
	}
	if got.BaseWorktreePath != want.BaseWorktreePath {
		t.Fatalf("repo base path mismatch: got %q want %q", got.BaseWorktreePath, want.BaseWorktreePath)
	}
}

func TestEnsureRepoRegisteredForRunRejectsInvalidName(t *testing.T) {
	fogHome := t.TempDir()
	store, err := state.NewStore(fogHome)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	_, err = ensureRepoRegisteredForRun("../repo", store, fogHome)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestEnsureRepoRegisteredForRunErrorsWhenGhMissing(t *testing.T) {
	fogHome := t.TempDir()
	store, err := state.NewStore(fogHome)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	origAvail := isGhAvailableFn
	t.Cleanup(func() { isGhAvailableFn = origAvail })
	isGhAvailableFn = func() bool { return false }

	_, err = ensureRepoRegisteredForRun("acme/api", store, fogHome)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "gh CLI not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnsureRepoRegisteredForRunErrorsWhenNotAuthenticated(t *testing.T) {
	fogHome := t.TempDir()
	store, err := state.NewStore(fogHome)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	origAvail := isGhAvailableFn
	origAuth := isGhAuthenticatedFn
	t.Cleanup(func() {
		isGhAvailableFn = origAvail
		isGhAuthenticatedFn = origAuth
	})
	isGhAvailableFn = func() bool { return true }
	isGhAuthenticatedFn = func() bool { return false }

	_, err = ensureRepoRegisteredForRun("acme/api", store, fogHome)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not authenticated") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnsureRepoRegisteredForRunErrorsWhenNotAccessible(t *testing.T) {
	fogHome := t.TempDir()
	store, err := state.NewStore(fogHome)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	origAvail := isGhAvailableFn
	origAuth := isGhAuthenticatedFn
	origList := listGitHubReposFn
	t.Cleanup(func() {
		isGhAvailableFn = origAvail
		isGhAuthenticatedFn = origAuth
		listGitHubReposFn = origList
	})
	isGhAvailableFn = func() bool { return true }
	isGhAuthenticatedFn = func() bool { return true }

	listGitHubReposFn = func() ([]ghcli.Repo, error) {
		return []ghcli.Repo{
			{NameWithOwner: "acme/other", Name: "other"},
		}, nil
	}

	_, err = ensureRepoRegisteredForRun("acme/api", store, fogHome)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "is not accessible") {
		t.Fatalf("expected not accessible error, got %v", err)
	}
}

func TestEnsureRepoRegisteredForRunAutoImports(t *testing.T) {
	fogHome := t.TempDir()
	store, err := state.NewStore(fogHome)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	origAvail := isGhAvailableFn
	origAuth := isGhAuthenticatedFn
	origList := listGitHubReposFn
	origClone := cloneGhRepoFn
	origRunner := gitRunner
	t.Cleanup(func() {
		isGhAvailableFn = origAvail
		isGhAuthenticatedFn = origAuth
		listGitHubReposFn = origList
		cloneGhRepoFn = origClone
		gitRunner = origRunner
	})

	isGhAvailableFn = func() bool { return true }
	isGhAuthenticatedFn = func() bool { return true }

	listGitHubReposFn = func() ([]ghcli.Repo, error) {
		r := ghcli.Repo{
			NameWithOwner: "acme/api",
			Name:          "api",
			URL:           "https://github.com/acme/api",
		}
		r.DefaultBranchRef.Name = "main"
		return []ghcli.Repo{r}, nil
	}

	cloneGhRepoFn = func(fullName, destPath string) error {
		if fullName != "acme/api" {
			return fmt.Errorf("unexpected clone name: %q", fullName)
		}
		return os.MkdirAll(destPath, 0o755)
	}

	gitRunner = func(extraEnv []string, args ...string) error {
		_ = extraEnv
		if len(args) > 0 && args[0] == "--git-dir" {
			basePath := args[len(args)-1]
			return os.MkdirAll(basePath, 0o755)
		}
		return fmt.Errorf("unexpected git args: %v", args)
	}

	repo, err := ensureRepoRegisteredForRun("acme/api", store, fogHome)
	if err != nil {
		t.Fatalf("ensureRepoRegisteredForRun returned error: %v", err)
	}
	if repo.Name != "acme/api" {
		t.Fatalf("repo name mismatch: got %q want %q", repo.Name, "acme/api")
	}
	wantBase := filepath.Join(fogHome, "repos", "acme", "api", "base")
	if repo.BaseWorktreePath != wantBase {
		t.Fatalf("base path mismatch: got %q want %q", repo.BaseWorktreePath, wantBase)
	}
}
