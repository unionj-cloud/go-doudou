package fileutils

import "os"

// If dir didn't exist, then create dir, otherwise do nothing.
func CreateDirectory(dir string) (err error) {
	if _, err = os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return
			}
		}
	}
	return
}
