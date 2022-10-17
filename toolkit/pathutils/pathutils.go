package pathutils

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"os"
	"path/filepath"
	"runtime"
)

// Abs converts relative path to absolute path
func Abs(rel string) string {
	_, fileName, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(fileName), rel)
}

// FixPath fixes path
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
