package codegen

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var mainTmpl = `package main

import (
	"context"
	"{{.RouterPackage}}"
	"flag"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "grace", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router.NewRouter(), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logrus.Infof("server is listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			logrus.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logrus.Infoln("shutting down")

	os.Exit(0)
}
`

func GenMain(dir string) {
	var (
		err       error
		modfile   string
		modName   string
		mainfile  string
		firstLine string
		f         *os.File
		tpl       *template.Template
		cmdDir    string
	)
	cmdDir = filepath.Join(dir, "cmd")
	if err = os.MkdirAll(cmdDir, os.ModePerm); err != nil {
		panic(err)
	}

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
			RouterPackage string
		}{
			RouterPackage: modName + "/router",
		}); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", mainfile)
	}
}
