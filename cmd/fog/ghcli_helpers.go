package main

import "github.com/darkLord19/foglet/internal/ghcli"

var (
	isGhAvailableFn     = ghcli.IsGhAvailable
	isGhAuthenticatedFn = ghcli.IsGhAuthenticated
	discoverGhReposFn   = ghcli.DiscoverRepos
	cloneGhRepoFn       = ghcli.CloneRepo
)
