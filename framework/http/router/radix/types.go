package radix

import (
	"regexp"

	"github.com/valyala/fasthttp"
)

type nodeType uint8

type nodeWildcard struct {
	path     string
	paramKey string
	handler  fasthttp.RequestHandler
}

type node struct {
	nType nodeType

	path         string
	tsr          bool
	handler      fasthttp.RequestHandler
	hasWildChild bool
	children     []*node
	wildcard     *nodeWildcard

	paramKeys  []string
	paramRegex *regexp.Regexp
}

type wildPath struct {
	path  string
	keys  []string
	start int
	end   int
	pType nodeType

	pattern string
	regex   *regexp.Regexp
}

// Tree is a routes storage
type Tree struct {
	root *node

	// If enabled, the node handler could be updated
	Mutable bool
}
