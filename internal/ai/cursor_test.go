package ai

import (
	"reflect"
	"testing"
)

func TestBuildCursorHeadlessArgsWithOutputFormat(t *testing.T) {
	got := buildCursorHeadlessArgs(ExecuteRequest{
		Prompt:         "fix auth",
		Model:          "gpt-5",
		ConversationID: "cursor-session-1",
	}, true)
	want := []string{"-p", "--force", "--model", "gpt-5", "--resume", "cursor-session-1", "--output-format", "stream-json", "fix auth"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args mismatch: got %v want %v", got, want)
	}
}

func TestBuildCursorHeadlessArgsWithoutOutputFormat(t *testing.T) {
	got := buildCursorHeadlessArgs(ExecuteRequest{Prompt: "fix auth"}, false)
	want := []string{"-p", "--force", "fix auth"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args mismatch: got %v want %v", got, want)
	}
}
