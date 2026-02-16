package main

import (
	"fmt"
	"regexp"
	"strings"
)

var repoSegmentPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

func splitRepoFullName(fullName string) (owner string, name string, err error) {
	fullName = strings.TrimSpace(fullName)
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("repo must be in owner/repo format")
	}
	owner = strings.TrimSpace(parts[0])
	name = strings.TrimSpace(parts[1])
	if owner == "" || name == "" {
		return "", "", fmt.Errorf("repo owner and name are required")
	}
	if owner == "." || owner == ".." || name == "." || name == ".." {
		return "", "", fmt.Errorf("repo contains invalid segment")
	}
	if !repoSegmentPattern.MatchString(owner) || !repoSegmentPattern.MatchString(name) {
		return "", "", fmt.Errorf("repo contains invalid characters")
	}
	return owner, name, nil
}
