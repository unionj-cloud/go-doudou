package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"
	"github.com/unionj-cloud/toolkit/pathutils"
	v3 "github.com/unionj-cloud/toolkit/protobuf/v3"
	"github.com/unionj-cloud/toolkit/stringutils"
)

var dbDriver string
var dbDsn string
var dbOrm string
var dbSoft string
var dbService bool
var dbTablePrefix string
var dbTableGlob string
var dbTableExcludeGlob string
var dbGenGenGo bool
var dbOmitempty bool

var crudCmd = &cobra.Command{
	Use:   "crud",
	Short: "generate universal crud code from database",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if stringutils.IsEmpty(dbDriver) || stringutils.IsEmpty(dbDsn) {
			logrus.Warn("Parameters db_driver and db_dsn must be specified")
			return
		}
		var options []svc.SvcOption
		options = append(options, svc.WithDbConfig(&svc.DbConfig{
			Driver:           dbDriver,
			Dsn:              dbDsn,
			TablePrefix:      dbTablePrefix,
			TableGlob:        dbTableGlob,
			TableExcludeGlob: dbTableExcludeGlob,
			GenGenGo:         dbGenGenGo,
			Orm:              dbOrm,
			Soft:             dbSoft,
			Service:          dbService,
			Omitempty:        dbOmitempty,
		}))
		fn := strcase.ToLowerCamel
		switch naming {
		case "snake":
			fn = strcase.ToSnake
		}
		options = append(options, svc.WithJsonCase(naming), svc.WithCaseConverter(fn), svc.WithProtoGenerator(v3.NewProtoGenerator(v3.WithFieldNamingFunc(fn), v3.WithProtocCmd(protocCmd))))
		svcdir, _ := pathutils.FixPath("", "")
		s := svc.NewSvc(svcdir, options...)
		s.Crud()
	},
}

func init() {
	svcCmd.AddCommand(crudCmd)

	crudCmd.Flags().StringVar(&dbOrm, "db_orm", "gorm", `Specify your preferable orm, currently only support gorm`)
	crudCmd.Flags().StringVar(&dbDriver, "db_driver", "", `Choose one database driver from "mysql", "postgres", "sqlite", "sqlserver", "tidb"`)
	crudCmd.Flags().StringVar(&dbDsn, "db_dsn", "", `Specify database connection url`)
	crudCmd.Flags().StringVar(&dbSoft, "db_soft", "deleted_at", `Specify database soft delete column name`)
	crudCmd.Flags().BoolVar(&dbService, "db_service", false, `If false, only dao layer code will be generated.`)
	crudCmd.Flags().BoolVar(&dbGenGenGo, "db_gen_gen", false, `whether generate gen.go file`)
	crudCmd.Flags().StringVar(&dbTablePrefix, "db_table_prefix", "", `table prefix or schema name for pg`)
	crudCmd.Flags().StringVar(&dbTableGlob, "db_table_glob", "", `used to filter glob-matched tables`)
	crudCmd.Flags().StringVar(&dbTableExcludeGlob, "db_table_exclude_glob", "", `used to filter glob-matched tables`)
	crudCmd.Flags().StringVar(&naming, "case", "lowerCamel", `protobuf message field and json tag case, only support "lowerCamel" and "snake"`)
	crudCmd.Flags().StringVar(&protocCmd, "grpc_gen_cmd", "protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --go-json_out=. --go-json_opt=paths=source_relative", `command to generate grpc service and message code`)
	crudCmd.Flags().BoolVar(&dbOmitempty, "db_omitempty", false, `whether add omitempty json tag to generated model field"`)
}
