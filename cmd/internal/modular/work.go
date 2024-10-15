package modular

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/toolkit/common"
	"github.com/unionj-cloud/toolkit/executils"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

const (
	mainDirName = "main"
	cmdDirName  = "cmd"
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

const workTmpl = `go 1.22.2

toolchain go1.22.3
`

func (receiver *Work) goModInMainPkg(dir, modName, goVersion string) {
	modfile := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(modfile); os.IsNotExist(err) {
		f, err := os.Create(modfile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ := template.New(templates.MainPkgModTmpl).Parse(templates.MainPkgModTmpl)
		_ = tpl.Execute(f, struct {
			ModName         string
			GoVersion       string
			GoDoudouVersion string
		}{
			ModName:         modName,
			GoVersion:       goVersion,
			GoDoudouVersion: version.Release,
		})
	} else {
		logrus.Warnf("file %s already exists", modfile)
	}
}

const envTmpl = ``

func (receiver *Work) dotenv(dir string) {
	envfile := filepath.Join(dir, ".env")
	if _, err := os.Stat(envfile); os.IsNotExist(err) {
		f, err := os.Create(envfile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ := template.New(envTmpl).Parse(envTmpl)
		_ = tpl.Execute(f, struct{}{})
	} else {
		logrus.Warnf("file %s already exists", envfile)
	}
}

func (receiver *Work) mainGo(dir string) {
	mainGoFile := filepath.Join(dir, "main.go")
	if _, err := os.Stat(mainGoFile); os.IsNotExist(err) {
		f, err := os.Create(mainGoFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		tpl, _ := template.New(templates.MainMainTmpl).Parse(templates.MainMainTmpl)
		_ = tpl.Execute(f, struct {
			Version string
		}{
			Version: version.Release,
		})
	} else {
		logrus.Warnf("file %s already exists", mainGoFile)
	}
}

func (receiver *Work) cmdPkg(dir string) {
	cmdDir := filepath.Join(dir, cmdDirName)
	if _, err := os.Stat(cmdDir); os.IsNotExist(err) {
		err = os.MkdirAll(cmdDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
		receiver.mainGo(cmdDir)
	} else {
		logrus.Warnf("directory %s already exists", cmdDir)
	}
}

func (receiver *Work) mainPkg(goVersion string) {
	mainDir := filepath.Join(receiver.GetWorkDir(), mainDirName)
	if _, err := os.Stat(mainDir); os.IsNotExist(err) {
		err = os.MkdirAll(mainDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
		modName := filepath.Base(receiver.GetWorkDir()) + "/" + mainDirName
		receiver.goModInMainPkg(mainDir, modName, goVersion)
		receiver.dotenv(mainDir)
		receiver.cmdPkg(mainDir)
	} else {
		logrus.Warnf("directory %s already exists", mainDir)
	}
}

func (receiver *Work) Init() {
	if stringutils.IsEmpty(receiver.GetWorkDir()) {
		wd, _ := os.Getwd()
		receiver.SetWorkDir(wd)
	}
	workDir := receiver.GetWorkDir()
	_ = os.MkdirAll(workDir, os.ModePerm)

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

	receiver.mainPkg(goVersion)

	if err = os.Chdir(receiver.GetWorkDir()); err != nil {
		panic(err)
	}
	if err = receiver.runner.Run("go", "work", "use", "main"); err != nil {
		panic(err)
	}
	// Comment below code due to performance issue
	//if err = receiver.runner.Run("go", "work", "sync"); err != nil {
	//	panic(err)
	//}
}
