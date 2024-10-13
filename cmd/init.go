package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"
	"github.com/unionj-cloud/toolkit/pathutils"
	v3 "github.com/unionj-cloud/toolkit/protobuf/v3"
)

var modName string
var module bool

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
		fn := strcase.ToLowerCamel
		switch naming {
		case "snake":
			fn = strcase.ToSnake
		}
		options = append(options, svc.WithJsonCase(naming), svc.WithCaseConverter(fn), svc.WithProtoGenerator(v3.NewProtoGenerator(v3.WithFieldNamingFunc(fn), v3.WithProtocCmd(protocCmd))))
		s := svc.NewSvc(svcdir, options...)
		s.Init()
	},
}

func init() {
	svcCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVar(&module, "module", false, `If true, a module will be initialized for building modular application`)
	initCmd.Flags().StringVarP(&modName, "mod", "m", "", `Module name`)
	initCmd.Flags().StringVarP(&docfile, "file", "f", "", `OpenAPI 3.0 or Swagger 2.0 spec json file path or download link`)
	initCmd.Flags().StringVar(&protocCmd, "grpc_gen_cmd", "protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --go-json_out=. --go-json_opt=paths=source_relative", `command to generate grpc service and message code`)
	initCmd.Flags().StringVar(&naming, "case", "lowerCamel", `protobuf message field and json tag case, only support "lowerCamel" and "snake"`)
}
