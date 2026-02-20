package runner

import (
	"fmt"
	"strings"
)

func resolvePRTitle(customTitle, prompt string) string {
	if title := firstNonEmptyLine(customTitle); title != "" {
		return truncate(title, 256)
	}
	if line := firstNonEmptyLine(prompt); line != "" {
		return fmt.Sprintf("feat: %s", truncate(line, 120))
	}
	return "feat: update code"
}

func firstNonEmptyLine(text string) string {
	for line := range strings.SplitSeq(text, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
