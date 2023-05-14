package stringutils

import (
	"regexp"
	"strings"
	"unicode"
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

func ToTitle(s string) string {
	str := []rune(s)
	return strings.ToUpper(string(str[0])) + string(str[1:])
}

var symbolre = regexp.MustCompile("[。？！，、；：“”‘’（）《》【】~!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~\\s]")

func ToCamel(s string) string {
	parts := symbolre.Split(s, -1)
	var convertedParts []string
	for _, v := range parts {
		if IsNotEmpty(v) {
			convertedParts = append(convertedParts, ToTitle(v))
		}
	}
	ret := strings.Join(convertedParts, "")
	if !unicode.IsUpper(rune(ret[0])) {
		ret = "A" + ret
	}
	return ret
}
