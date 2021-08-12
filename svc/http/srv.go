package ddhttp

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"github.com/unionj-cloud/go-doudou/svc/http/model"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	logFile *os.File
)

func init() {
	configureLogger()
}

type Srv interface {
	// Run the service
	Run()
	// Register routes
	AddRoute(route ...model.Route)
	// Use middleware
	AddMiddleware(mwf ...func(http.Handler) http.Handler)
}

func newServer(router http.Handler) *http.Server {
	write, err := time.ParseDuration(config.GddWriteTimeout.Load())
	if err != nil {
		logrus.Warnf("Parse %s %s as time.Duration failed: %s, use default 15s instead.\n", "GDD_WRITETIMEOUT",
			config.GddWriteTimeout.Load(), err.Error())
		write = 15 * time.Second
	}

	read, err := time.ParseDuration(config.GddReadTimeout.Load())
	if err != nil {
		logrus.Warnf("Parse %s %s as time.Duration failed: %s, use default 15s instead.\n", "GDD_READTIMEOUT",
			config.GddReadTimeout.Load(), err.Error())
		read = 15 * time.Second
	}

	idle, err := time.ParseDuration(config.GddIdleTimeout.Load())
	if err != nil {
		logrus.Warnf("Parse %s %s as time.Duration failed: %s, use default 60s instead.\n", "GDD_IDLETIMEOUT",
			config.GddIdleTimeout.Load(), err.Error())
		idle = 60 * time.Second
	}

	server := &http.Server{
		Addr: strings.Join([]string{config.GddHost.Load(), config.GddPort.Load()}, ":"),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: write,
		ReadTimeout:  read,
		IdleTimeout:  idle,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logrus.Infof("Http server is listening on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			logrus.Println(err)
		}
	}()

	return server
}

func configureLogger() {
	var logptr *string
	logpath, isSet := os.LookupEnv(config.GddLogPath.String())
	if isSet {
		logptr = &logpath
	}
	var loglevel config.LogLevel
	(&loglevel).Decode(config.GddLogLevel.Load())

	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true

	logger := logrus.StandardLogger()
	logger.SetFormatter(formatter)
	logger.SetLevel(logrus.Level(loglevel))

	if logptr != nil {
		var err error
		logpath := *logptr
		logpath, err = pathutils.FixPath(logpath, "")
		if err != nil {
			logger.Errorln(fmt.Sprintf("%+v\n", err))
		}
		if stringutils.IsNotEmpty(logpath) {
			err = os.MkdirAll(logpath, os.ModePerm)
			if err != nil {
				logger.Errorln(err)
			}
		}
		logFile, err = os.OpenFile(filepath.Join(logpath, "app.log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			logger.Errorf("error opening file: %v\n", err)
		}
		mw := io.MultiWriter(os.Stdout, logFile)
		logger.SetOutput(mw)
	}
}

func printRoutes(routes []model.Route) {
	logrus.Infoln("================ Registered Routes ================")
	data := [][]string{}
	for _, r := range routes {
		data = append(data, []string{r.Name, r.Method, config.GddRouteRootPath.Load() + r.Pattern})
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Name", "Method", "Pattern"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
	rows := strings.Split(strings.TrimSpace(tableString.String()), "\n")
	for _, row := range rows {
		logrus.Infoln(row)
	}
	logrus.Infoln("===================================================")
}
