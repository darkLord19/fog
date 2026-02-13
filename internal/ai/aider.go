package ai

import (
	"context"
	"fmt"
	"strings"
)

// Aider represents the Aider AI tool
type Aider struct{}

func (a *Aider) Name() string {
	return "aider"
}

func (a *Aider) IsAvailable() bool {
	return commandExists("aider")
}

func (a *Aider) Execute(ctx context.Context, workdir, prompt string) (*Result, error) {
	return a.ExecuteStream(ctx, ExecuteRequest{
		Workdir: workdir,
		Prompt:  prompt,
	}, nil)
}

func (a *Aider) ExecuteStream(ctx context.Context, req ExecuteRequest, onChunk func(string)) (*Result, error) {
	if !a.IsAvailable() {
		return nil, fmt.Errorf("aider not available")
	}
	cmdName := commandPath("aider")
	if cmdName == "" {
		return nil, fmt.Errorf("aider not available")
	}

	args := []string{"--yes"}
	if model := strings.TrimSpace(req.Model); model != "" {
		args = append(args, "--model", model)
	}
	args = append(args, "--message", strings.TrimSpace(req.Prompt))

	output, err := runPlainStreamingCommand(ctx, req.Workdir, cmdName, args, onChunk)

	result := &Result{
		Success: err == nil,
		Output:  strings.TrimSpace(output),
	}
	if err != nil {
		result.Error = err
	}
	return result, err
}
