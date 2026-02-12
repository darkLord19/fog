package slack

import (
	"fmt"
	"regexp"
	"strings"
)

var optionPattern = regexp.MustCompile(`([a-zA-Z-]+)=('([^']*)'|"([^"]*)"|[^\s]+)`)

type parsedCommand struct {
	Repo       string
	Tool       string
	Model      string
	AutoPR     bool
	BranchName string
	CommitMsg  string
	Prompt     string
}

func parseCommandText(raw string) (*parsedCommand, error) {
	text := strings.TrimSpace(raw)
	if strings.HasPrefix(strings.ToLower(text), "@fog") {
		text = strings.TrimSpace(text[len("@fog"):])
	}
	if text == "" {
		return nil, fmt.Errorf("invalid command format. Use: @fog [repo='' tool='' model='' autopr=true/false branch-name='' commit-msg=''] prompt")
	}

	if !strings.HasPrefix(text, "[") {
		return nil, fmt.Errorf("options block is required. Use: @fog [repo='...'] prompt")
	}
	end := strings.Index(text, "]")
	if end == -1 {
		return nil, fmt.Errorf("invalid options block: missing closing ]")
	}

	optsText := strings.TrimSpace(text[1:end])
	prompt := strings.TrimSpace(text[end+1:])
	if prompt == "" {
		return nil, fmt.Errorf("prompt is required")
	}

	opts, err := parseOptions(optsText)
	if err != nil {
		return nil, err
	}

	repo := strings.TrimSpace(opts["repo"])
	if repo == "" {
		return nil, fmt.Errorf("repo is required")
	}

	autopr := false
	if rawVal, ok := opts["autopr"]; ok && strings.TrimSpace(rawVal) != "" {
		val := strings.ToLower(strings.TrimSpace(rawVal))
		switch val {
		case "true":
			autopr = true
		case "false":
			autopr = false
		default:
			return nil, fmt.Errorf("invalid autopr value %q, expected true/false", rawVal)
		}
	}

	return &parsedCommand{
		Repo:       repo,
		Tool:       strings.TrimSpace(opts["tool"]),
		Model:      strings.TrimSpace(opts["model"]),
		AutoPR:     autopr,
		BranchName: strings.TrimSpace(opts["branch-name"]),
		CommitMsg:  strings.TrimSpace(opts["commit-msg"]),
		Prompt:     prompt,
	}, nil
}

func parseOptions(input string) (map[string]string, error) {
	allowed := map[string]struct{}{
		"repo":        {},
		"tool":        {},
		"model":       {},
		"autopr":      {},
		"branch-name": {},
		"commit-msg":  {},
	}

	matches := optionPattern.FindAllStringSubmatchIndex(input, -1)
	if len(matches) == 0 && strings.TrimSpace(input) != "" {
		return nil, fmt.Errorf("invalid options format")
	}

	options := make(map[string]string)
	cursor := 0
	for _, m := range matches {
		if strings.TrimSpace(input[cursor:m[0]]) != "" {
			return nil, fmt.Errorf("invalid options format near %q", strings.TrimSpace(input[cursor:m[0]]))
		}

		key := strings.ToLower(input[m[2]:m[3]])
		if _, ok := allowed[key]; !ok {
			return nil, fmt.Errorf("unknown option key: %s", key)
		}

		value := input[m[4]:m[5]]
		value = strings.TrimPrefix(value, "'")
		value = strings.TrimSuffix(value, "'")
		value = strings.TrimPrefix(value, "\"")
		value = strings.TrimSuffix(value, "\"")
		options[key] = value
		cursor = m[1]
	}

	if strings.TrimSpace(input[cursor:]) != "" {
		return nil, fmt.Errorf("invalid options format near %q", strings.TrimSpace(input[cursor:]))
	}

	return options, nil
}

func generateBranchName(prefix, prompt string) string {
	cleanPrefix := sanitizeSegment(prefix)
	if cleanPrefix == "" {
		cleanPrefix = "fog"
	}

	slug := sanitizeSegment(prompt)
	if slug == "" {
		slug = "task"
	}

	branch := cleanPrefix + "/" + slug
	if len(branch) <= 255 {
		return branch
	}

	maxSlug := 255 - len(cleanPrefix) - 1
	if maxSlug < 1 {
		return cleanPrefix[:255]
	}
	if len(slug) > maxSlug {
		slug = slug[:maxSlug]
	}
	return cleanPrefix + "/" + strings.Trim(slug, "-/")
}

func sanitizeSegment(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(input))
	lastDash := false

	for _, r := range input {
		isAlphaNum := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		switch {
		case isAlphaNum:
			b.WriteRune(r)
			lastDash = false
		case r == '/' || r == '-' || r == '_' || r == '.':
			if !lastDash {
				b.WriteRune(r)
				lastDash = true
			}
		default:
			if !lastDash {
				b.WriteRune('-')
				lastDash = true
			}
		}
	}

	out := strings.Trim(b.String(), "-/_.")
	return out
}

func isProtectedBranch(branch string) bool {
	name := strings.ToLower(strings.TrimSpace(branch))
	return name == "main" || name == "master"
}
