package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
)

var modName string
var dbDriver string
var dbDsn string
var dbOrm string
var dbSoft string

// initCmd initializes the service
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a project folder",
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
		options := []svc.SvcOption{svc.WithModName(modName), svc.WithDocPath(docfile)}
		dbConf := svc.DbConfig{
			Driver: dbDriver,
			Dsn:    dbDsn,
			Orm:    dbOrm,
			Soft:   dbSoft,
		}
		if stringutils.IsNotEmpty(dbConf.Driver) && stringutils.IsNotEmpty(dbConf.Dsn) {
			options = append(options, svc.WithDbConfig(&dbConf))
		}
		s := svc.NewSvc(svcdir, options...)
		s.Init()
	},
}

func init() {
	svcCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&modName, "mod", "m", "", `Module name`)
	initCmd.Flags().StringVarP(&docfile, "file", "f", "", `OpenAPI 3.0 or Swagger 2.0 spec json file path or download link`)
	initCmd.Flags().StringVar(&dbOrm, "db_orm", "gorm", `Specify your preferable orm, currently only support gorm`)
	initCmd.Flags().StringVar(&dbDriver, "db_driver", "", `Choose one database driver from "mysql", "postgres", "sqlite", "sqlserver", "tidb"`)
	initCmd.Flags().StringVar(&dbDsn, "db_dsn", "", `Specify database connection url`)
	initCmd.Flags().StringVar(&dbSoft, "db_soft", "deleted_at", `Specify database soft delete column name`)
}
