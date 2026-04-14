package audit

import "unicode/utf8"

// TruncateQuery stores at most maxRunes runes of s for audit metadata (help avoid huge payloads).
func TruncateQuery(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes]) + "…"
}
