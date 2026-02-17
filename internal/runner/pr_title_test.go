package runner

import "testing"

func TestResolvePRTitlePrefersCustomTitle(t *testing.T) {
	got := resolvePRTitle("\n  feat: custom title  \n\nmore details\n", "Add JWT auth")
	if got != "feat: custom title" {
		t.Fatalf("unexpected title: got %q want %q", got, "feat: custom title")
	}
}

func TestResolvePRTitleFallsBackToPromptFirstLine(t *testing.T) {
	got := resolvePRTitle("", "Add JWT auth\n\nMore details here")
	if got != "feat: Add JWT auth" {
		t.Fatalf("unexpected title: got %q want %q", got, "feat: Add JWT auth")
	}
}

func TestResolvePRTitleHandlesEmptyInput(t *testing.T) {
	got := resolvePRTitle("", "")
	if got != "feat: update code" {
		t.Fatalf("unexpected title: got %q want %q", got, "feat: update code")
	}
}

func TestFirstNonEmptyLineReturnsEmptyForWhitespace(t *testing.T) {
	got := firstNonEmptyLine(" \n\t\n")
	if got != "" {
		t.Fatalf("unexpected line: got %q want empty", got)
	}
}
