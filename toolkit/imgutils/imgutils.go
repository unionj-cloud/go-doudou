package imgutils

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/toolkit/caller"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

func ResizeKeepAspectRatio(input io.Reader, multiplier float64, output string) (string, error) {
	by, err := io.ReadAll(input)
	if err != nil {
		return "", errors.Wrap(err, caller.NewCaller().String())
	}
	imgConf, imgType, err := image.DecodeConfig(bytes.NewReader(by))
	if err != nil {
		return "", errors.Wrap(err, caller.NewCaller().String())
	}
	switch imgType {
	case "jpeg":
		img, err := jpeg.Decode(bytes.NewReader(by))
		if err != nil {
			return "", errors.Wrap(err, caller.NewCaller().String())
		}
		m := resize.Resize(uint(float64(imgConf.Width)*multiplier), 0, img, resize.Lanczos3)
		if stringutils.IsEmpty(filepath.Ext(output)) {
			output += fmt.Sprintf(".%s", imgType)
		}
		out, err := os.Create(output)
		if err != nil {
			return "", errors.Wrap(err, caller.NewCaller().String())
		}
		defer out.Close()
		jpeg.Encode(out, m, nil)
	case "png":
		img, err := png.Decode(bytes.NewReader(by))
		if err != nil {
			return "", errors.Wrap(err, caller.NewCaller().String())
		}
		m := resize.Resize(uint(float64(imgConf.Width)*multiplier), 0, img, resize.Lanczos3)
		if stringutils.IsEmpty(filepath.Ext(output)) {
			output += fmt.Sprintf(".%s", imgType)
		}
		out, err := os.Create(output)
		if err != nil {
			return "", errors.Wrap(err, caller.NewCaller().String())
		}
		defer out.Close()
		png.Encode(out, m)
	case "gif":
		img, err := gif.DecodeAll(bytes.NewReader(by))
		if err != nil {
			return "", errors.Wrap(err, caller.NewCaller().String())
		}
		if multiplier != 1 {
			for i, item := range img.Image {
				resized := resize.Resize(uint(float64(imgConf.Width)*multiplier), 0, item, resize.Lanczos3)
				rgba64 := resized.(*image.RGBA64)
				palettedImage := image.NewPaletted(rgba64.Bounds(), getSubPalette(rgba64))
				draw.Draw(palettedImage, palettedImage.Rect, rgba64, rgba64.Bounds().Min, draw.Src)
				img.Image[i] = palettedImage
			}
		}
		if stringutils.IsEmpty(filepath.Ext(output)) {
			output += fmt.Sprintf(".%s", imgType)
		}
		out, err := os.Create(output)
		if err != nil {
			return "", errors.Wrap(err, caller.NewCaller().String())
		}
		defer out.Close()
		gif.EncodeAll(out, img)
	}
	return output, nil
}
