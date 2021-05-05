package codegen

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var mainTmpl = `package main

import (
	"context"
	{{.ServiceAlias}} "{{.ServicePackage}}"
	"{{.ConfigPackage}}"
	"{{.DbPackage}}"
	"{{.HttpPackage}}"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func configureLogger(logger *logrus.Logger, logptr *string, level logrus.Level) *os.File {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	logger.SetFormatter(formatter)
	logger.SetLevel(level)

	if logptr != nil {
		var (
			err error
			f   *os.File
		)
		logpath := *logptr
		logpath, err = pathutils.FixPath(logpath, "")
		if err != nil {
			logger.Errorln(fmt.Sprintf("%+v\n", err))
		}
		if stringutils.IsNotEmpty(logpath) {
			err = os.MkdirAll(logpath, os.ModePerm)
			if err != nil {
				logger.Errorln(err)
				return nil
			}
		}
		f, err = os.OpenFile(filepath.Join(logpath, "app.log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			logger.Errorf("error opening file: %v\n", err)
			return nil
		}
		mw := io.MultiWriter(os.Stdout, f)
		logger.SetOutput(mw)
		return f
	}

	return nil
}

func main() {
	start := time.Now()
	env := config.NewDotenv(pathutils.Abs("../.env"))
	conf := env.Get()

	logFile := configureLogger(logrus.StandardLogger(), conf.AppConf.LogPath, logrus.Level(conf.AppConf.LogLevel))
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	if conf.AppConf.Banner {
		figure.NewColorFigure(conf.AppConf.BannerText, "doom", "green", true).Print()
	}

	db, err := db.NewDb(conf.DbConf)
	if err != nil {
		panic(err)
	}
	defer func() {
		if db == nil {
			return
		}
		if err := db.Close(); err == nil {
			logrus.Infoln("Database connection is closed")
		} else {
			logrus.Warnln("Failed to close database connection")
		}
	}()

	svc := {{.ServiceAlias}}.New{{.SvcName}}(conf.SvcConf, db)
	handler := httpsrv.New{{.SvcName}}Handler(svc)

	srv := httpsrv.Run(conf.HttpConf, handler)

	logrus.Infof("Started in %s\n", time.Since(start))

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), conf.AppConf.GraceTimeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logrus.Infoln("shutting down")
}
`

func GenMain(dir string, ic astutils.InterfaceCollector) {
	var (
		err       error
		modfile   string
		modName   string
		mainfile  string
		firstLine string
		f         *os.File
		tpl       *template.Template
		cmdDir    string
		svcName   string
		alias     string
	)
	cmdDir = filepath.Join(dir, "cmd")
	if err = os.MkdirAll(cmdDir, os.ModePerm); err != nil {
		panic(err)
	}

	svcName = ic.Interfaces[0].Name
	alias = ic.Package.Name
	mainfile = filepath.Join(cmdDir, "main.go")
	if _, err = os.Stat(mainfile); os.IsNotExist(err) {
		modfile = filepath.Join(dir, "go.mod")
		if f, err = os.Open(modfile); err != nil {
			panic(err)
		}
		reader := bufio.NewReader(f)
		if firstLine, err = reader.ReadString('\n'); err != nil {
			panic(err)
		}
		modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))

		if f, err = os.Create(mainfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("main.go.tmpl").Parse(mainTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, struct {
			ServicePackage string
			ConfigPackage  string
			DbPackage      string
			HttpPackage    string
			SvcName        string
			ServiceAlias   string
		}{
			ServicePackage: modName,
			ConfigPackage:  modName + "/config",
			DbPackage:      modName + "/db",
			HttpPackage:    modName + "/transport/httpsrv",
			SvcName:        svcName,
			ServiceAlias:   alias,
		}); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", mainfile)
	}
}
