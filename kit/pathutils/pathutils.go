package pathutils

import (
	"path/filepath"
	"runtime"
)

func Abs(rel string) string {
	_, fileName, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(fileName), rel)
}