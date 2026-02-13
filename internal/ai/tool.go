package ai

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Tool represents an AI coding tool
type Tool interface {
	Name() string
	IsAvailable() bool
	Execute(ctx context.Context, workdir, prompt string) (*Result, error)
}

// StreamTool is implemented by tools that support incremental output streaming.
type StreamTool interface {
	Tool
	ExecuteStream(ctx context.Context, req ExecuteRequest, onChunk func(string)) (*Result, error)
}

// ExecuteRequest configures one tool execution call.
type ExecuteRequest struct {
	Workdir        string
	Prompt         string
	Model          string
	ConversationID string
}

// Result contains the AI execution result
type Result struct {
	Success        bool
	Output         string
	Error          error
	ConversationID string
}

// GetTool returns an AI tool by name
func GetTool(name string) (Tool, error) {
	switch normalizeToolName(name) {
	case "cursor":
		return &Cursor{}, nil
	case "claude", "claude-code":
		return &ClaudeCode{}, nil
	case "gemini":
		return &Gemini{}, nil
	case "aider":
		return &Aider{}, nil
	default:
		return nil, fmt.Errorf("unknown AI tool: %s", name)
	}
}

// DetectTool finds an available AI tool
func DetectTool(preferred string) (Tool, error) {
	tools := []Tool{
		&Cursor{},
		&ClaudeCode{},
		&Gemini{},
		&Aider{},
	}

	// Try preferred first
	if preferred != "" {
		tool, err := GetTool(preferred)
		if err == nil && tool.IsAvailable() {
			return tool, nil
		}
	}

	// Fall back to first available
	for _, tool := range tools {
		if tool.IsAvailable() {
			return tool, nil
		}
	}

	return nil, fmt.Errorf("no AI tool available")
}

// AvailableToolNames returns canonical tool names supported by Fog.
func AvailableToolNames() []string {
	return []string{"cursor", "claude", "gemini", "aider"}
}

// ExecuteWithOptionalStream runs a tool and streams chunks when supported.
func ExecuteWithOptionalStream(ctx context.Context, tool Tool, req ExecuteRequest, onChunk func(string)) (*Result, error) {
	if streamTool, ok := tool.(StreamTool); ok {
		return streamTool.ExecuteStream(ctx, req, onChunk)
	}
	result, err := tool.Execute(ctx, req.Workdir, req.Prompt)
	if onChunk != nil && result != nil && strings.TrimSpace(result.Output) != "" {
		onChunk(result.Output)
	}
	return result, err
}

func normalizeToolName(name string) string {
	value := strings.ToLower(strings.TrimSpace(name))
	switch value {
	case "claude-code":
		return "claude"
	default:
		return value
	}
}

// commandExists checks if a command is available
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
