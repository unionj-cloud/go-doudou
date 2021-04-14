package pathutils

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"os"
	"path/filepath"
	"runtime"
)

func Abs(rel string) string {
	_, fileName, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(fileName), rel)
}

func FixPath(dir string, fallback string) (string, error) {
	var (
		wd  string
		err error
		ret string
	)
	if stringutils.IsEmpty(dir) {
		if wd, err = os.Getwd(); err != nil {
			return ret, errors.Wrap(err, "Error from calling os.Getwd()")
		}
		ret = filepath.Join(wd, fallback)
		return ret, nil
	}
	if !filepath.IsAbs(dir) {
		if wd, err = os.Getwd(); err != nil {
			logrus.Panicln(err)
		}
		ret = filepath.Join(wd, dir)
		return ret, nil
	}
	return dir, nil
}
