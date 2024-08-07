package common

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
)

// InitGitRepo inits git repository.
// Reinitialized existing Git repository is safe
// https://stackoverflow.com/questions/5149694/does-running-git-init-twice-initialize-a-repository-or-reinitialize-an-existing
func InitGitRepo(dir string) {
	fs := osfs.New(dir)
	dot, _ := fs.Chroot(".git")
	storage := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	_, _ = git.Init(storage, fs)
}

const gitignoreTmpl = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/
**/*.local
.DS_Store
.idea`

// GitIgnore adds .gitignore file
func GitIgnore(dir string) {
	var (
		gitignorefile string
		err           error
		f             *os.File
		tpl           *template.Template
	)
	gitignorefile = filepath.Join(dir, ".gitignore")
	if _, err = os.Stat(gitignorefile); os.IsNotExist(err) {
		if f, err = os.Create(gitignorefile); err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ = template.New(".gitignore.tmpl").Parse(gitignoreTmpl)
		_ = tpl.Execute(f, nil)
	} else {
		logrus.Warnf("file %s already exists", ".gitignore")
	}
}

func GetGoVersionNum(runner executils.Runner) (string, error) {
	out, err := runner.Output("go", "version")
	if err != nil {
		return "", errors.WithStack(err)
	}
	// go version go1.18.8 darwin/amd64
	result := strings.TrimPrefix(strings.Split(string(out), " ")[2], "go")
	// invalid go version '1.18.8': must match format 1.23
	return result[:strings.LastIndex(result, ".")], nil
}
