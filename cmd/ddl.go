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
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/ddl"
	"github.com/unionj-cloud/go-doudou/pathutils"

	"github.com/spf13/cobra"
)

var dir string
var reverse bool
var dao bool
var pre string
var df string
var env string

// ddlCmd represents the ddl command
var ddlCmd = &cobra.Command{
	Use:   "ddl",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if env, err = pathutils.FixPath(env, ".env"); err != nil {
			logrus.Panicln(err)
		}
		if err = godotenv.Load(env); err != nil {
			logrus.Panicln("Error loading .env file", err)
		}
		if dir, err = pathutils.FixPath(dir, "domain"); err != nil {
			logrus.Panicln(err)
		}

		d := ddl.Ddl{dir, reverse, dao, pre, df}
		d.Exec()
	},
}

func init() {
	rootCmd.AddCommand(ddlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ddlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	ddlCmd.Flags().StringVar(&dir, "domain", "domain", "Path of domain folder.")
	ddlCmd.Flags().StringVar(&pre, "pre", "", "Table name prefix. e.g.: prefix biz_ for biz_product.")
	ddlCmd.Flags().StringVar(&df, "df", "dao", "Name of dao folder.")
	ddlCmd.Flags().StringVar(&env, "env", ".env", "Path of database connection config .env file")
	ddlCmd.Flags().BoolVarP(&reverse, "reverse", "r", false, "If true, generate domain code from database. If false, update or create database tables from domain code.")
	ddlCmd.Flags().BoolVarP(&dao, "dao", "d", false, "If true, generate dao code.")
}
