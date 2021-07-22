/*
Copyright © 2021 wubin1989 <328454505@qq.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/svc"

	"github.com/spf13/cobra"
)

var docfile string
var lang string
var baseUrlEnv string
var clientpkg string

// clientCmd represents the http command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "generate http client from openapi 3.0 spec json file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var svcdir string
		if len(args) > 0 {
			svcdir = args[0]
		}
		var err error
		if svcdir, err = pathutils.FixPath(svcdir, ""); err != nil {
			logrus.Panicln(err)
		}
		s := svc.Svc{
			Dir:       svcdir,
			DocPath:   docfile,
			Client:    lang,
			Omitempty: omitempty,
			Env:       baseUrlEnv,
			ClientPkg: clientpkg,
		}
		s.GenClient()
	},
}

func init() {
	httpCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&lang, "lang", "l", "go", `client language`)
	clientCmd.Flags().StringVarP(&docfile, "file", "f", "", `openapi 3.0 spec json file path or download link`)
	clientCmd.Flags().StringVarP(&baseUrlEnv, "env", "e", "", `base url environment variable name`)
	clientCmd.Flags().StringVarP(&clientpkg, "pkg", "p", "client", `client package name`)
	clientCmd.Flags().BoolVarP(&omitempty, "omit", "o", false, `json tag omitempty`)
}
