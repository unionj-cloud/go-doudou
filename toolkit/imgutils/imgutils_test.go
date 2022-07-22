package imgutils_test

import (
	"github.com/unionj-cloud/go-doudou/toolkit/imgutils"
	"os"
	"path/filepath"
	"testing"
)

const testDir = "testdata"

func TestResizeKeepAspectRatioPng(t *testing.T) {
	file, err := os.Open(filepath.Join(testDir, "test.png"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = imgutils.ResizeKeepAspectRatio(file, 0.5, filepath.Join(testDir, "test_result"))
	if err != nil {
		panic(err)
	}
}

func TestResizeKeepAspectRatioJpeg(t *testing.T) {
	file, err := os.Open(filepath.Join(testDir, "test.jpg"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = imgutils.ResizeKeepAspectRatio(file, 0.5, filepath.Join(testDir, "test_result.jpg"))
	if err != nil {
		panic(err)
	}
}

func TestResizeKeepAspectRatioGif1(t *testing.T) {
	file, err := os.Open(filepath.Join(testDir, "rgb.gif"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = imgutils.ResizeKeepAspectRatio(file, 0.5, filepath.Join(testDir, "rgb_result"))
	if err != nil {
		panic(err)
	}
}

func TestResizeKeepAspectRatioGif2(t *testing.T) {
	file, err := os.Open(filepath.Join(testDir, "test.gif"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = imgutils.ResizeKeepAspectRatio(file, 1, filepath.Join(testDir, "test_result"))
	if err != nil {
		panic(err)
	}
}

func TestResizeKeepAspectRatioGif3(t *testing.T) {
	file, err := os.Open(filepath.Join(testDir, "test.gif"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = imgutils.ResizeKeepAspectRatio(file, 0.5, filepath.Join(testDir, "test_result"))
	if err != nil {
		panic(err)
	}
}
