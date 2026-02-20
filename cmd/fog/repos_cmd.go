package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	fogenv "github.com/darkLord19/foglet/internal/env"
	"github.com/darkLord19/foglet/internal/ghcli"
	"github.com/darkLord19/foglet/internal/state"
	"github.com/spf13/cobra"
)

var (
	reposJSONFlag   bool
	reposSelectFlag string
	gitRunner       = runGitCommand
)

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Manage Fog repositories",
}

var reposDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "List repositories accessible by GitHub CLI",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runReposDiscover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var reposImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Select and register repositories from GitHub using gh CLI",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runReposImport(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var reposListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories already registered in Fog",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runReposList(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	reposDiscoverCmd.Flags().BoolVar(&reposJSONFlag, "json", false, "Output JSON")
	reposImportCmd.Flags().StringVar(&reposSelectFlag, "select", "", "Comma-separated GitHub full names to import (e.g. org/repo,org/repo2)")

	reposCmd.AddCommand(reposDiscoverCmd)
	reposCmd.AddCommand(reposImportCmd)
	reposCmd.AddCommand(reposListCmd)
	rootCmd.AddCommand(reposCmd)
}

func runReposDiscover() error {
	repos, err := discoverGitHubRepos()
	if err != nil {
		return err
	}

	if reposJSONFlag {
		data, err := json.MarshalIndent(repos, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	if len(repos) == 0 {
		fmt.Println("No accessible repositories found via gh CLI")
		return nil
	}

	fmt.Printf("%-40s %-8s %s\n", "FULL NAME", "PRIVATE", "DEFAULT BRANCH")
	fmt.Println(strings.Repeat("-", 70))
	for _, repo := range repos {
		fmt.Printf("%-40s %-8t %s\n", repo.NameWithOwner, repo.IsPrivate, repo.DefaultBranchRef.Name)
	}

	return nil
}

func runReposImport() error {
	repos, err := discoverGitHubRepos()
	if err != nil {
		return err
	}
	if len(repos) == 0 {
		fmt.Println("No repositories available to import")
		return nil
	}

	selected, err := selectRepos(repos, reposSelectFlag)
	if err != nil {
		return err
	}
	if len(selected) == 0 {
		fmt.Println("No repositories selected")
		return nil
	}

	fogHome, err := fogenv.FogHome()
	if err != nil {
		return err
	}

	store, err := state.NewStore(fogHome)
	if err != nil {
		return err
	}
	defer func() { _ = store.Close() }()

	managedReposDir := fogenv.ManagedReposDir(fogHome)
	if err := os.MkdirAll(managedReposDir, 0o755); err != nil {
		return fmt.Errorf("create managed repos dir: %w", err)
	}

	for _, repo := range selected {
		owner, name, err := splitRepoFullName(repo.NameWithOwner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping invalid repo name %q: %v\n", repo.NameWithOwner, err)
			continue
		}

		repoDir := filepath.Join(managedReposDir, owner, name)
		if err := os.MkdirAll(repoDir, 0o755); err != nil {
			return fmt.Errorf("create repo dir %s: %w", repoDir, err)
		}
		barePath := filepath.Join(repoDir, "repo.git")
		basePath := filepath.Join(repoDir, "base")

		if err := ensureBareRepoInitialized(repo, barePath, basePath); err != nil {
			return err
		}

		host := repoHost(repo.URL)
		_, err = store.UpsertRepo(state.Repo{
			Name:             repo.NameWithOwner,
			URL:              repo.URL,
			Host:             host,
			Owner:            owner,
			Repo:             name,
			BarePath:         barePath,
			BaseWorktreePath: basePath,
			DefaultBranch:    repo.DefaultBranchRef.Name,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Imported %s\n", repo.NameWithOwner)
	}

	return nil
}

func runReposList() error {
	fogHome, err := fogenv.FogHome()
	if err != nil {
		return err
	}
	store, err := state.NewStore(fogHome)
	if err != nil {
		return err
	}
	defer func() { _ = store.Close() }()

	repos, err := store.ListRepos()
	if err != nil {
		return err
	}
	if len(repos) == 0 {
		fmt.Println("No repositories registered")
		return nil
	}

	fmt.Printf("%-40s %-40s %s\n", "NAME", "URL", "DEFAULT BRANCH")
	fmt.Println(strings.Repeat("-", 100))
	for _, repo := range repos {
		fmt.Printf("%-40s %-40s %s\n", repo.Name, repo.URL, repo.DefaultBranch)
	}

	return nil
}

func discoverGitHubRepos() ([]ghcli.Repo, error) {
	if !isGhAvailableFn() {
		return nil, fmt.Errorf("gh CLI invalid or not found")
	}
	if !isGhAuthenticatedFn() {
		return nil, fmt.Errorf("gh CLI not authenticated; run `gh auth login`")
	}
	return discoverGhReposFn()
}

func ensureBareRepoInitialized(repo ghcli.Repo, barePath, basePath string) error {
	if _, err := os.Stat(barePath); errorsIsNotExist(err) {
		if err := cloneGhRepoFn(repo.NameWithOwner, barePath); err != nil {
			return fmt.Errorf("clone bare repository %s: %w", repo.NameWithOwner, err)
		}
	} else if err != nil {
		return fmt.Errorf("check bare repo path %s: %w", barePath, err)
	}

	if _, err := os.Stat(basePath); errorsIsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(basePath), 0o755); err != nil {
			return fmt.Errorf("create base worktree parent: %w", err)
		}
		if err := gitRunner(nil, "--git-dir", barePath, "worktree", "add", basePath); err != nil {
			return fmt.Errorf("create base worktree for %s: %w", repo.NameWithOwner, err)
		}
	} else if err != nil {
		return fmt.Errorf("check base worktree path %s: %w", basePath, err)
	}

	return nil
}

func runGitCommand(extraEnv []string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Env = append(os.Environ(), extraEnv...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return nil
}

func errorsIsNotExist(err error) bool {
	return err != nil && os.IsNotExist(err)
}

func selectRepos(repos []ghcli.Repo, selectFlag string) ([]ghcli.Repo, error) {
	selectFlag = strings.TrimSpace(selectFlag)
	if selectFlag == "" {
		fmt.Println("Available repositories:")
		for i, repo := range repos {
			fmt.Printf("  %d. %s\n", i+1, repo.NameWithOwner)
		}
		input, err := readLine("Select repository numbers (comma-separated): ")
		if err != nil {
			return nil, err
		}
		indexes, err := parseIndexes(input, len(repos))
		if err != nil {
			return nil, err
		}
		selected := make([]ghcli.Repo, 0, len(indexes))
		for _, idx := range indexes {
			selected = append(selected, repos[idx])
		}
		return selected, nil
	}

	want := make(map[string]struct{})
	for name := range strings.SplitSeq(selectFlag, ",") {
		name = strings.TrimSpace(name)
		if name != "" {
			want[name] = struct{}{}
		}
	}
	if len(want) == 0 {
		return nil, fmt.Errorf("--select cannot be empty")
	}

	selected := make([]ghcli.Repo, 0, len(want))
	for _, repo := range repos {
		if _, ok := want[repo.NameWithOwner]; ok {
			selected = append(selected, repo)
			delete(want, repo.NameWithOwner)
		}
	}
	if len(want) > 0 {
		missing := make([]string, 0, len(want))
		for name := range want {
			missing = append(missing, name)
		}
		return nil, fmt.Errorf("unknown repositories in --select: %s", strings.Join(missing, ", "))
	}

	return selected, nil
}

func parseIndex(input string, max int) (int, error) {
	i, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, fmt.Errorf("invalid selection %q", input)
	}
	if i < 1 || i > max {
		return 0, fmt.Errorf("selection out of range: %d", i)
	}
	return i - 1, nil
}

func parseIndexes(input string, max int) ([]int, error) {
	parts := strings.Split(input, ",")
	result := make([]int, 0, len(parts))
	seen := make(map[int]struct{})
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idx, err := parseIndex(part, max)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[idx]; ok {
			continue
		}
		seen[idx] = struct{}{}
		result = append(result, idx)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no valid selections provided")
	}
	return result, nil
}

func repoAlias(repo ghcli.Repo) string {
	return repo.NameWithOwner
}

func repoHost(cloneURL string) string {
	u, err := url.Parse(cloneURL)
	if err != nil {
		return "github.com"
	}
	if u.Host == "" {
		return "github.com"
	}
	return u.Host
}
