package branch

import (
	"regexp"
	"strings"
)

func Generate(branchType, ticket string, description []string) string {
	var parts []string

	if ticket != "" {
		parts = append(parts, ticket)
	}

	parts = append(parts, description...)

	slug := slugify(strings.Join(parts, " "))

	if slug == "" {
		return branchType
	}

	return branchType + "/" + slug
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)

	// Replace spaces and underscores with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove any characters that aren't alphanumeric, hyphens, or forward slashes
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	s = reg.ReplaceAllString(s, "")

	// Collapse multiple hyphens
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	// Trim leading/trailing hyphens
	s = strings.Trim(s, "-")

	return s
}
