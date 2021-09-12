package pathutils

import (
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
		ret string
	)
	if stringutils.IsEmpty(dir) {
		wd, _ = os.Getwd()
		ret = filepath.Join(wd, fallback)
		return ret, nil
	}
	if !filepath.IsAbs(dir) {
		wd, _ = os.Getwd()
		ret = filepath.Join(wd, dir)
		return ret, nil
	}
	return dir, nil
}
