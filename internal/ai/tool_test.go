package ai

import "slices"

import "testing"

func TestGetToolGemini(t *testing.T) {
	tool, err := GetTool("gemini")
	if err != nil {
		t.Fatalf("GetTool returned error: %v", err)
	}
	if tool.Name() != "gemini" {
		t.Fatalf("unexpected tool name: got %q want %q", tool.Name(), "gemini")
	}
}

func TestGetToolClaudeAlias(t *testing.T) {
	tool, err := GetTool("claude-code")
	if err != nil {
		t.Fatalf("GetTool returned error: %v", err)
	}
	if tool.Name() != "claude" {
		t.Fatalf("unexpected tool name: got %q want %q", tool.Name(), "claude")
	}
}

func TestAvailableToolNamesIncludesGemini(t *testing.T) {
	names := AvailableToolNames()
	found := slices.Contains(names, "gemini")
	if !found {
		t.Fatalf("expected gemini in available tool names: %v", names)
	}
}
