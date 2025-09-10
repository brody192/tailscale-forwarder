package util

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"syscall"
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

func IsExpectedCopyError(err error) bool {
	if err == nil {
		return true
	}

	if errors.Is(err, io.EOF) ||
		errors.Is(err, net.ErrClosed) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.EPIPE) ||
		errors.Is(err, context.Canceled) {
		return true
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}

	return false
}
