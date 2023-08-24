package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/protobuf/v3"

	"github.com/spf13/cobra"
)

var naming string

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "generate grpc service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fn := strcase.ToLowerCamel
		switch naming {
		case "snake":
			fn = strcase.ToSnake
		}
		s := svc.NewSvc("", svc.WithProtoGenerator(v3.NewProtoGenerator(v3.WithFieldNamingFunc(fn))))
		s.Grpc()
	},
}

func init() {
	svcCmd.AddCommand(grpcCmd)
	grpcCmd.Flags().StringVar(&naming, "case", "lowerCamel", `protobuf message field naming strategy, only support "lowerCamel" and "snake"`)
}
