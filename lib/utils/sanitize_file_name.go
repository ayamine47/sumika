package utils

import (
	"regexp"
	"strings"
)

var invalidCharsRegex = regexp.MustCompile(`[/\\:*?"<>|]`)

func SanitizeFileName(name string) string {
	cleaned := invalidCharsRegex.ReplaceAllString(name, "_")
	return strings.TrimSpace(cleaned)
}
