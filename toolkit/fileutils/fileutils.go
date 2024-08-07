package fileutils

import (
	"bufio"
	"io"
	"os"

	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
)

// CreateDirectory dir didn't exist, then create dir, otherwise do nothing.
func CreateDirectory(dir string) (err error) {
	var info os.FileInfo
	if info, err = os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return
			}
		}
	} else {
		if !info.IsDir() {
			return errors.New("not a directory: " + dir)
		}
	}
	return
}

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func Archive(output string, sources ...string) error {
	return archiver.Archive(sources, output)
}
