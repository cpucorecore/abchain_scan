package util

import "unicode/utf8"

// TruncateToMaxChars safely truncates a string to a maximum number of characters
// without breaking UTF-8 encoding. If the string exceeds maxChars,
// it returns a substring containing at most maxChars complete runes.
func TruncateToMaxChars(s string, maxChars int) string {
	if utf8.RuneCountInString(s) <= maxChars {
		return s
	}
	// Convert to runes for proper character handling
	runes := []rune(s)
	return string(runes[:maxChars])
}
