package gormgen

import "strings"

var (
	importList = new(importPkgS).Add(
		"context",
		"database/sql",
		"strings",
		"",
		"github.com/wubin1989/gorm",
		"github.com/wubin1989/gorm/schema",
		"github.com/wubin1989/gorm/clause",
		"",
		"github.com/unionj-cloud/go-doudou/v2/toolkit/gormgen",
		"github.com/unionj-cloud/go-doudou/v2/toolkit/gormgen/field",
		"github.com/unionj-cloud/go-doudou/v2/toolkit/gormgen/helper",
		"",
		"github.com/wubin1989/dbresolver",
	)
	unitTestImportList = new(importPkgS).Add(
		"context",
		"fmt",
		"strconv",
		"testing",
		"",
		"github.com/wubin1989/sqlite",
		"github.com/wubin1989/gorm",
	)
)

type importPkgS struct {
	paths []string
}

func (ip importPkgS) Add(paths ...string) *importPkgS {
	purePaths := make([]string, 0, len(paths)+1)
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			purePaths = append(purePaths, p)
			continue
		}

		if p[len(p)-1] != '"' {
			p = `"` + p + `"`
		}

		var exists bool
		for _, existsP := range ip.paths {
			if p == existsP {
				exists = true
				break
			}
		}
		if !exists {
			purePaths = append(purePaths, p)
		}
	}
	purePaths = append(purePaths, "")

	ip.paths = append(ip.paths, purePaths...)

	return &ip
}

func (ip importPkgS) Paths() []string { return ip.paths }
