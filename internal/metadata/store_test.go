package metadata

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewUsesGitCommonDirInRepoRoot(t *testing.T) {
	repo := initTestRepo(t)

	store, err := New(repo)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	want := filepath.Join(repo, ".git", "wtx", "metadata.json")
	if canonicalPath(t, store.path) != canonicalPath(t, want) {
		t.Fatalf("unexpected metadata path: got %q want %q", store.path, want)
	}
}

func TestNewUsesGitCommonDirInWorktree(t *testing.T) {
	repo := initTestRepo(t)
	worktreePath := filepath.Join(filepath.Dir(repo), "feature-wt")

	runGit(t, repo, "worktree", "add", "-b", "feature", worktreePath)

	store, err := New(worktreePath)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	want := filepath.Join(repo, ".git", "wtx", "metadata.json")
	if canonicalPath(t, store.path) != canonicalPath(t, want) {
		t.Fatalf("unexpected metadata path for worktree: got %q want %q", store.path, want)
	}
}

func initTestRepo(t *testing.T) string {
	t.Helper()

	repo := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatalf("mkdir repo failed: %v", err)
	}

	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "fog-test@example.com")
	runGit(t, repo, "config", "user.name", "fog test")

	filePath := filepath.Join(repo, "README.md")
	if err := os.WriteFile(filePath, []byte("test\n"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	runGit(t, repo, "add", "README.md")
	runGit(t, repo, "commit", "-m", "init")

	return repo
}

func runGit(t *testing.T, repo string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(out))
	}
}

func canonicalPath(t *testing.T, path string) string {
	t.Helper()
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	resolvedDir, err := filepath.EvalSymlinks(dir)
	if err == nil {
		return filepath.Join(resolvedDir, base)
	}
	return filepath.Clean(path)
}
