package timeutils

import (
	"github.com/hyperjumptech/jiffy"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"time"
)

// Parse parses string to time.Duration
func Parse(t string, defaultDur time.Duration) (time.Duration, error) {
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
