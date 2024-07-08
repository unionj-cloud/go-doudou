package stringutils

import (
	"regexp"
	"strings"
	"unicode"
)

// IsEmpty asserts s is empty
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == "" || strings.TrimSpace(s) == "<nil>"
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

func ReplaceAtRuneIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func ReplaceStringAtByteIndex(in string, replace string, start int, end int) string {
	out := []byte(in)
	r := []byte(replace)
	result := make([]byte, len(out[:start]))
	copy(result, out[:start])
	result = append(result, r...)
	result = append(result, out[end:]...)
	return string(result)
}

func ReplaceStringAtByteIndexBatch(in string, args []string, locs [][]int) string {
	out := []byte(in)
	result := make([]byte, 0)
	end := 0
	for i, loc := range locs {
		arg := args[i]
		r := []byte(arg)
		start := loc[0]
		result = append(result, out[end:start]...)
		result = append(result, r...)
		end = loc[1]
	}
	result = append(result, out[end:]...)
	return string(result)
}
