package main

import (
	"reflect"
	"testing"

	foggithub "github.com/darkLord19/wtx/pkg/fog/github"
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

func TestRepoAlias(t *testing.T) {
	repo := foggithub.Repo{FullName: "acme/api", OwnerLogin: "acme", Name: "api"}
	got := repoAlias(repo)
	if got != "acme-api" {
		t.Fatalf("repoAlias mismatch: got %q want %q", got, "acme-api")
	}
}

func TestSelectReposByFullName(t *testing.T) {
	repos := []foggithub.Repo{
		{FullName: "acme/api", Name: "api", OwnerLogin: "acme"},
		{FullName: "acme/web", Name: "web", OwnerLogin: "acme"},
	}
	selected, err := selectRepos(repos, "acme/web")
	if err != nil {
		t.Fatalf("selectRepos returned error: %v", err)
	}
	if len(selected) != 1 || selected[0].FullName != "acme/web" {
		t.Fatalf("unexpected selected repos: %+v", selected)
	}
}
