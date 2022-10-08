package radix

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/valyala/bytebufferpool"
)

func panicf(s string, args ...interface{}) {
	panic(fmt.Sprintf(s, args...))
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func bufferRemoveString(buf *bytebufferpool.ByteBuffer, s string) {
	buf.B = buf.B[:len(buf.B)-len(s)]
}

// func isIndexEqual(a, b string) bool {
// 	ra, _ := utf8.DecodeRuneInString(a)
// 	rb, _ := utf8.DecodeRuneInString(b)

// 	return unicode.ToLower(ra) == unicode.ToLower(rb)
// }

// longestCommonPrefix finds the longest common prefix.
// This also implies that the common prefix contains no ':' or '*'
// since the existing key can't contain those chars.
func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(utf8.RuneCountInString(a), utf8.RuneCountInString(b))

	for i < max {
		ra, sizeA := utf8.DecodeRuneInString(a)
		rb, sizeB := utf8.DecodeRuneInString(b)

		a = a[sizeA:]
		b = b[sizeB:]

		if ra != rb {
			return i
		}

		i += sizeA
	}

	return i
}

// segmentEndIndex returns the index where the segment ends from the given path
func segmentEndIndex(path string, includeTSR bool) int {
	end := 0
	for end < len(path) && path[end] != '/' {
		end++
	}

	if includeTSR && path[end:] == "/" {
		end++
	}

	return end
}

// findWildPath search for a wild path segment and check the name for invalid characters.
// Returns -1 as index, if no param/wildcard was found.
func findWildPath(path string, fullPath string) *wildPath {
	// Find start
	for start, c := range []byte(path) {
		// A wildcard starts with ':' (param) or '*' (wildcard)
		if c != '{' {
			continue
		}

		withRegex := false
		keys := 0

		// Find end and check for invalid characters
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '}':
				if keys > 0 {
					keys--
					continue
				}

				end := start + end + 2
				wp := &wildPath{
					path:  path[start:end],
					keys:  []string{path[start+1 : end-1]},
					start: start,
					end:   end,
					pType: param,
				}

				if len(path) > end && path[end] == '{' {
					panic("the wildcards must be separated by at least 1 char")
				}

				sn := strings.SplitN(wp.keys[0], ":", 2)
				if len(sn) > 1 {
					wp.keys = []string{sn[0]}
					pattern := sn[1]

					if pattern == "*" {
						wp.pattern = pattern
						wp.pType = wildcard
					} else {
						wp.pattern = "(" + pattern + ")"
						wp.regex = regexp.MustCompile(wp.pattern)
					}
				} else if path[len(path)-1] != '/' {
					wp.pattern = "(.*)"
				}

				if len(wp.keys[0]) == 0 {
					panicf("wildcards must be named with a non-empty name in path '%s'", fullPath)
				}

				segEnd := end + segmentEndIndex(path[end:], true)
				path = path[end:segEnd]

				if path == "/" {
					// Last segment, so include the TSR
					path = ""
					wp.end++
				}

				if len(path) > 0 {
					// Rebuild the wildpath with the prefix
					wp2 := findWildPath(path, fullPath)
					if wp2 != nil {
						prefix := path[:wp2.start]

						wp.end += wp2.end
						wp.path += prefix + wp2.path
						wp.pattern += prefix + wp2.pattern
						wp.keys = append(wp.keys, wp2.keys...)
					} else {
						wp.path += path
						wp.pattern += path
						wp.end += len(path)
					}

					wp.regex = regexp.MustCompile(wp.pattern)
				}

				return wp

			case ':':
				withRegex = true

			case '{':
				if !withRegex && keys == 0 {
					panic("the char '{' is not allowed in the param name")
				}

				keys++
			}
		}
	}

	return nil
}
