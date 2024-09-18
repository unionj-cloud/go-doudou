package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/protobuf/v3"

	"github.com/spf13/cobra"
)

var naming string
var http2grpc bool
var annotatedOnly bool

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
		s := svc.NewSvc("",
			svc.WithProtoGenerator(v3.NewProtoGenerator(v3.WithFieldNamingFunc(fn), v3.WithAnnotatedOnly(annotatedOnly))),
			svc.WithHttp2Grpc(http2grpc),
			svc.WithAllowGetWithReqBody(allowGetWithReqBody),
			svc.WithCaseConverter(fn),
			svc.WithOmitempty(omitempty),
		)
		s.Grpc()
	},
}

func init() {
	svcCmd.AddCommand(grpcCmd)
	grpcCmd.Flags().BoolVarP(&omitempty, "omitempty", "o", false, `if true, ",omitempty" will be appended to json tag of fields in every generated anonymous struct in handlers`)
	grpcCmd.Flags().StringVar(&naming, "case", "lowerCamel", `protobuf message field naming strategy, only support "lowerCamel" and "snake"`)
	grpcCmd.Flags().BoolVar(&http2grpc, "http2grpc", false, `whether need RESTful api for your grpc service`)
	grpcCmd.Flags().BoolVar(&allowGetWithReqBody, "allowGetWithReqBody", false, "Whether allow get http request with request body.")
	grpcCmd.Flags().BoolVar(&annotatedOnly, "annotatedOnly", false, "Whether generate grpc api only for method annotated with @grpc or not")
}
