package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBranchExistsAndAddWorktreeNewBranch(t *testing.T) {
	repo := initGitRepo(t)
	g := New(repo)

	if !g.BranchExists("main") {
		t.Fatal("expected main branch to exist")
	}
	if g.BranchExists("feature") {
		t.Fatal("did not expect feature branch to exist yet")
	}

	wtPath := filepath.Join(filepath.Dir(repo), "feature-wt")
	if err := g.AddWorktreeNewBranch(wtPath, "feature", "main"); err != nil {
		t.Fatalf("AddWorktreeNewBranch failed: %v", err)
	}

	if !g.BranchExists("feature") {
		t.Fatal("expected feature branch to exist after worktree creation")
	}
	if _, err := os.Stat(wtPath); err != nil {
		t.Fatalf("expected worktree path to exist: %v", err)
	}
}

func initGitRepo(t *testing.T) string {
	t.Helper()

	repo := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatalf("mkdir repo failed: %v", err)
	}

	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "fog-test@example.com")
	runGit(t, repo, "config", "user.name", "fog test")

	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
	runGit(t, repo, "add", "README.md")
	runGit(t, repo, "commit", "-m", "init")
	runGit(t, repo, "branch", "-M", "main")

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
