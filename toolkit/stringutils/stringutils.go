package stringutils

import (
	"regexp"
	"strings"
)

// IsEmpty asserts s is empty
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty asserts s is not empty
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// ContainsI assert s contains substr ignore case
func ContainsI(s string, substr string) bool {
	re := regexp.MustCompile(`(?i)` + substr)
	return re.MatchString(s)
}

// HasPrefixI assert s has prefix prefix ignore case
func HasPrefixI(s, prefix string) bool {
	re := regexp.MustCompile(`(?i)^` + prefix)
	return re.MatchString(s)
}
