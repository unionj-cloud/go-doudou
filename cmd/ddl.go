package cmd

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/ddl"
	"github.com/unionj-cloud/go-doudou/ddl/config"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"os"
)

var dir string
var reverse bool
var dao bool
var pre string
var df string
var env string

// ddlCmd generates domain and dao layer source code from database tables and update tables from domain code
var ddlCmd = &cobra.Command{
	Use:   "ddl",
	Short: "migration tool between database table structure and golang struct",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if env, err = pathutils.FixPath(env, ".env"); err != nil {
			logrus.Panicln(err)
		}
		if _, err = os.Stat(env); err == nil {
			if err = godotenv.Load(env); err != nil {
				logrus.Panicln("Error loading .env file", err)
			}
		}
		if dir, err = pathutils.FixPath(dir, "domain"); err != nil {
			logrus.Panicln(err)
		}
		var conf config.DbConfig
		err = envconfig.Process("db", &conf)
		if err != nil {
			logrus.Panicln("Error processing env", err)
		}
		d := ddl.Ddl{dir, reverse, dao, pre, df, conf}
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
