package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/protobuf/v3"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
)

var modName string
var dbDriver string
var dbDsn string
var dbOrm string
var dbSoft string
var dbGrpc bool
var dbService bool
var dbTablePrefix string
var dbTableGlob string
var module bool
var dbGenGenGo bool

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
		options := []svc.SvcOption{svc.WithModName(modName), svc.WithDocPath(docfile), svc.WithModule(module)}
		dbConf := svc.DbConfig{
			Driver:      dbDriver,
			Dsn:         dbDsn,
			TablePrefix: dbTablePrefix,
			TableGlob:   dbTableGlob,
			GenGenGo:    dbGenGenGo,
			Orm:         dbOrm,
			Soft:        dbSoft,
			Grpc:        dbGrpc,
			Service:     dbService,
		}
		if stringutils.IsNotEmpty(dbConf.Driver) && stringutils.IsNotEmpty(dbConf.Dsn) {
			options = append(options, svc.WithDbConfig(&dbConf))
		}
		fn := strcase.ToLowerCamel
		switch naming {
		case "snake":
			fn = strcase.ToSnake
		}
		options = append(options, svc.WithJsonCase(naming), svc.WithCaseConverter(fn), svc.WithProtoGenerator(v3.NewProtoGenerator(v3.WithFieldNamingFunc(fn))))
		s := svc.NewSvc(svcdir, options...)
		s.Init()
	},
}

func init() {
	svcCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVar(&module, "module", false, `If true, a module will be initialized for building modular application`)
	initCmd.Flags().StringVarP(&modName, "mod", "m", "", `Module name`)
	initCmd.Flags().StringVarP(&docfile, "file", "f", "", `OpenAPI 3.0 or Swagger 2.0 spec json file path or download link`)
	initCmd.Flags().StringVar(&dbOrm, "db_orm", "gorm", `Specify your preferable orm, currently only support gorm`)
	initCmd.Flags().StringVar(&dbDriver, "db_driver", "", `Choose one database driver from "mysql", "postgres", "sqlite", "sqlserver", "tidb"`)
	initCmd.Flags().StringVar(&dbDsn, "db_dsn", "", `Specify database connection url`)
	initCmd.Flags().StringVar(&dbSoft, "db_soft", "deleted_at", `Specify database soft delete column name`)
	initCmd.Flags().BoolVar(&dbGrpc, "db_grpc", false, `If true, grpc code will also be generated`)
	initCmd.Flags().BoolVar(&dbService, "db_service", false, `If false, service will not be generated, and db_grpc will be ignored. Only dao layer code will be generated.`)
	initCmd.Flags().BoolVar(&dbGenGenGo, "db_gen_gen", false, `whether generate gen.go file`)
	initCmd.Flags().StringVar(&dbTablePrefix, "db_table_prefix", "", `table prefix or schema name for pg`)
	initCmd.Flags().StringVar(&dbTableGlob, "db_table_glob", "", `used to filter glob-matched tables`)
	initCmd.Flags().StringVar(&naming, "case", "lowerCamel", `protobuf message field and json tag case, only support "lowerCamel" and "snake"`)
}
