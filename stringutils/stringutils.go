package stringutils

import (
	"regexp"
	"strings"
)

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// Case insensitive match
func ContainsI(s string, substr string) bool {
	re := regexp.MustCompile(`(?i)` + substr)
	return re.MatchString(s)
}

func HasPrefixI(s, prefix string) bool {
	re := regexp.MustCompile(`(?i)^` + prefix)
	return re.MatchString(s)
}
