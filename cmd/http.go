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

var handler bool
var client string
var homitempty bool
var doc bool
var jsonattrcase string

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
			Dir:          svcdir,
			Handler:      handler,
			Client:       client,
			Omitempty:    homitempty,
			Doc:          doc,
			Jsonattrcase: jsonattrcase,
		}
		s.Http()
	},
}

func init() {
	svcCmd.AddCommand(httpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// httpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	httpCmd.Flags().BoolVarP(&handler, "handler", "", false, "Whether generate default handler implementation or not")
	httpCmd.Flags().StringVarP(&client, "client", "c", "", `if empty, then no http client implementation will be generated. Only one value "go" supported currently`)
	httpCmd.Flags().BoolVarP(&homitempty, "omitempty", "o", false, `if true, ",omitempty" will be appended to json tag of fields in every generated anonymous struct in handlers`)
	httpCmd.Flags().StringVarP(&jsonattrcase, "case", "", "lowerCamel", `apply to json tag of fields in every generated anonymous struct in handlers. optional values: lowerCamel, snake`)
	httpCmd.Flags().BoolVarP(&doc, "doc", "", false, `whether generate openapi 3.0 json document or not`)
}
