package httpsrv

import (
	"example/user-svc/config"
	"fmt"
	"github.com/hyperjumptech/jiffy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"net/http"
	"strings"
	"time"
)

func parseTimeout(t string, defaultDur time.Duration) (time.Duration, error) {
	var (
		dur time.Duration
		err error
	)
	if stringutils.IsNotEmpty(t) {
		if dur, err = jiffy.DurationOf(t); err != nil {
			err = errors.Wrapf(err, "parse %s from config file fail, use default 15s instead", t)
		}
	}
	if dur <= 0 {
		dur = defaultDur
	}
	return dur, err
}

func Start(conf config.HttpConfig) *http.Server {
	host := conf.Host
	if stringutils.IsEmpty(host) {
		host = "0.0.0.0"
	}
	port := conf.Port
	if stringutils.IsEmpty(port) {
		port = "8080"
	}
	wt, err := parseTimeout(conf.WriteTimeout, time.Second*15)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("%+v", err))
	}

	rt, err := parseTimeout(conf.ReadTimeout, time.Second*15)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("%+v", err))
	}

	it, err := parseTimeout(conf.IdleTimeout, time.Second*60)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("%+v", err))
	}

	srv := &http.Server{
		Addr: strings.Join([]string{host, port}, ":"),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: wt,
		ReadTimeout:  rt,
		IdleTimeout:  it,
		Handler:      newRouter(), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logrus.Infof("server is listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			logrus.Println(err)
		}
	}()

	return srv
}
