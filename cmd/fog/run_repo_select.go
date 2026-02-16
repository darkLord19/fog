package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	fogenv "github.com/darkLord19/foglet/internal/env"
	"github.com/darkLord19/foglet/internal/ghcli"
	"github.com/darkLord19/foglet/internal/state"
	"golang.org/x/term"
)

var (
	listGitHubReposFn = listGitHubRepos
	readLineFn        = readLine
	stdinIsTTYFn      = stdinIsTTY
)

func resolveRepoNameForRun(flagRepo string, store *state.Store) (string, error) {
	flagRepo = strings.TrimSpace(flagRepo)
	if flagRepo != "" {
		owner, name, err := splitRepoFullName(flagRepo)
		if err != nil {
			return "", err
		}
		return owner + "/" + name, nil
	}

	if !stdinIsTTYFn() {
		return "", fmt.Errorf("--repo is required (owner/repo); run 'fog repos discover' to list accessible repos")
	}

	if !isGhAvailableFn() {
		return "", fmt.Errorf("gh CLI not found")
	}
	if !isGhAuthenticatedFn() {
		return "", fmt.Errorf("gh CLI not authenticated; run `gh auth login`")
	}

	repos, err := listGitHubReposFn()
	if err != nil {
		return "", err
	}
	if len(repos) == 0 {
		return "", fmt.Errorf("no accessible repositories found via gh CLI")
	}

	fmt.Println("Available repositories:")
	for i, repo := range repos {
		fmt.Printf("  %d. %s\n", i+1, repo.NameWithOwner)
	}

	input, err := readLineFn("Select repository number: ")
	if err != nil {
		return "", err
	}
	idx, err := parseIndex(input, len(repos))
	if err != nil {
		return "", err
	}

	name := strings.TrimSpace(repos[idx].NameWithOwner)
	if name == "" {
		return "", errors.New("selected repository has no name")
	}
	owner, repo, err := splitRepoFullName(name)
	if err != nil {
		return "", err
	}
	return owner + "/" + repo, nil
}

func ensureRepoRegisteredForRun(repoName string, store *state.Store, fogHome string) (state.Repo, error) {
	repoName = strings.TrimSpace(repoName)
	if repoName == "" {
		return state.Repo{}, fmt.Errorf("repo name is required")
	}

	owner, name, err := splitRepoFullName(repoName)
	if err != nil {
		return state.Repo{}, err
	}
	repoName = owner + "/" + name

	if repo, found, err := store.GetRepoByName(repoName); err != nil {
		return state.Repo{}, err
	} else if found {
		return repo, nil
	}

	if !isGhAvailableFn() {
		return state.Repo{}, fmt.Errorf("gh CLI not found")
	}
	if !isGhAuthenticatedFn() {
		return state.Repo{}, fmt.Errorf("gh CLI not authenticated; run `gh auth login`")
	}

	repos, err := listGitHubReposFn()
	if err != nil {
		return state.Repo{}, err
	}

	match, ok := findRepoByFullName(repos, repoName)
	if !ok {
		return state.Repo{}, fmt.Errorf("repo %q is not accessible via gh CLI", repoName)
	}

	owner, name, err = splitRepoFullName(match.NameWithOwner)
	if err != nil {
		return state.Repo{}, fmt.Errorf("invalid repo name %q: %v", match.NameWithOwner, err)
	}

	managedReposDir := fogenv.ManagedReposDir(fogHome)
	if err := os.MkdirAll(managedReposDir, 0o755); err != nil {
		return state.Repo{}, fmt.Errorf("create managed repos dir: %w", err)
	}

	repoDir := filepath.Join(managedReposDir, owner, name)
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		return state.Repo{}, fmt.Errorf("create repo dir %s: %w", repoDir, err)
	}
	barePath := filepath.Join(repoDir, "repo.git")
	basePath := filepath.Join(repoDir, "base")

	if err := ensureBareRepoInitialized(match, barePath, basePath); err != nil {
		return state.Repo{}, err
	}

	host := repoHost(match.URL)
	_, err = store.UpsertRepo(state.Repo{
		Name:             match.NameWithOwner,
		URL:              match.URL,
		Host:             host,
		Owner:            owner,
		Repo:             name,
		BarePath:         barePath,
		BaseWorktreePath: basePath,
		DefaultBranch:    match.DefaultBranchRef.Name,
	})
	if err != nil {
		return state.Repo{}, err
	}

	repo, found, err := store.GetRepoByName(match.NameWithOwner)
	if err != nil {
		return state.Repo{}, err
	}
	if !found {
		return state.Repo{}, fmt.Errorf("managed repo %q disappeared after import", match.NameWithOwner)
	}
	return repo, nil
}

func listGitHubRepos() ([]ghcli.Repo, error) {
	if !isGhAvailableFn() {
		return nil, fmt.Errorf("gh CLI invalid or not found")
	}
	if !isGhAuthenticatedFn() {
		return nil, fmt.Errorf("gh CLI not authenticated; run `gh auth login`")
	}
	return discoverGhReposFn()
}

func stdinIsTTY() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func findRepoByFullName(repos []ghcli.Repo, fullName string) (ghcli.Repo, bool) {
	fullName = strings.TrimSpace(fullName)
	for _, repo := range repos {
		if strings.TrimSpace(repo.NameWithOwner) == fullName {
			return repo, true
		}
	}
	return ghcli.Repo{}, false
}
