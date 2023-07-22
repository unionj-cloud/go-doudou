package modular

import (
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/common"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"os"
	"path/filepath"
	"text/template"
)

type WorkConfig struct {
	WorkDir string
}

type Work struct {
	conf   WorkConfig
	runner executils.Runner
}

func NewWork(conf WorkConfig, runner executils.Runner) *Work {
	return &Work{
		conf:   conf,
		runner: runner,
	}
}

func (receiver *Work) GetWorkDir() string {
	return receiver.conf.WorkDir
}

func (receiver *Work) SetWorkDir(workDir string) {
	receiver.conf.WorkDir = workDir
}

const workTmpl = `go {{.GoVersion}}
`

func (receiver *Work) Init() {
	if stringutils.IsEmpty(receiver.GetWorkDir()) {
		wd, _ := os.Getwd()
		receiver.SetWorkDir(wd)
	}
	workDir := receiver.GetWorkDir()
	_ = os.MkdirAll(workDir, os.ModePerm)

	common.InitGitRepo(workDir)
	common.GitIgnore(workDir)

	goVersion, err := common.GetGoVersionNum(receiver.runner)
	if err != nil {
		panic(err)
	}
	workFile := filepath.Join(receiver.GetWorkDir(), "go.work")
	if _, err = os.Stat(workFile); os.IsNotExist(err) {
		f, err := os.Create(workFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ := template.New(workTmpl).Parse(workTmpl)
		_ = tpl.Execute(f, struct {
			GoVersion string
		}{
			GoVersion: goVersion,
		})
	} else {
		logrus.Warnf("file %s already exists", workFile)
	}
}
