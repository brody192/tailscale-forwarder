package util

import (
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	var result strings.Builder

	for _, char := range input {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune('-')
		}
	}

	sanitized := result.String()
	for strings.Contains(sanitized, "--") {
		sanitized = strings.ReplaceAll(sanitized, "--", "-")
	}

	sanitized = strings.Trim(sanitized, "-")

	return sanitized
}
