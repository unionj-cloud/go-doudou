package cmd

import (
	"github.com/iancoleman/strcase"
	v3 "github.com/unionj-cloud/go-doudou/v2/cmd/internal/protobuf/v3"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"

	"github.com/spf13/cobra"
)

var naming string

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "generate grpc service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := svc.NewSvc("")
		fn := strcase.ToLowerCamel
		switch naming {
		case "snake":
			fn = strcase.ToSnake
		}
		s.Grpc(v3.NewProtoGenerator(v3.WithFieldNamingFunc(fn)))
	},
}

func init() {
	svcCmd.AddCommand(grpcCmd)
	grpcCmd.Flags().StringVarP(&naming, "naming", "n", "lowerCamel", `protobuf message field naming strategy, only support "lowerCamel" and "snake"`)
}
