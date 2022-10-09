package router

import (
	"strings"
)

// cleanPath removes the '.' if it is the last character of the route
func cleanPath(path string) string {
	return strings.TrimSuffix(path, ".")
}