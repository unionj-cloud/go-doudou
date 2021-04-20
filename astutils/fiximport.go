package astutils

import (
	"bytes"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
)

func FixImport(src []byte, file string) {
	var (
		res []byte
		err error
	)
	if res, err = imports.Process(file, src, &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}); err != nil {
		panic(err)
	}

	if !bytes.Equal(src, res) {
		// On Windows, we need to re-set the permissions from the file. See golang/go#38225.
		var perms os.FileMode
		var fi os.FileInfo
		if fi, err = os.Stat(file); err == nil {
			perms = fi.Mode() & os.ModePerm
		}
		err = ioutil.WriteFile(file, res, perms)
		if err != nil {
			panic(err)
		}
	}
}
