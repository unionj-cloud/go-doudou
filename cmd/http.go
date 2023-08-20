package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"

	"github.com/spf13/cobra"
)

var handler bool
var client bool
var doc bool
var jsonCase string
var routePatternStrategy int
var allowGetWithReqBody bool

// httpCmd generates scaffold code of restful service
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "generate http routes and handlers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fn := strcase.ToLowerCamel
		switch jsonCase {
		case "snake":
			fn = strcase.ToSnake
		}
		s := svc.Svc{
			Handler:              handler,
			Client:               client,
			Omitempty:            omitempty,
			Doc:                  doc,
			CaseConverter:        fn,
			Env:                  baseURLEnv,
			RoutePatternStrategy: routePatternStrategy,
			AllowGetWithReqBody:  allowGetWithReqBody,
		}
		s.Http()
	},
}

func init() {
	svcCmd.AddCommand(httpCmd)

	httpCmd.Flags().BoolVar(&handler, "handler", false, "Whether generate default handler implementation or not")
	httpCmd.Flags().BoolVarP(&client, "client", "c", false, `Whether generate default golang http client code or not`)
	httpCmd.Flags().BoolVarP(&omitempty, "omitempty", "o", false, `if true, ",omitempty" will be appended to json tag of fields in every generated anonymous struct in handlers`)
	httpCmd.Flags().StringVar(&jsonCase, "case", "lowerCamel", `apply to json tag of fields in every generated anonymous struct in handlers. optional values: lowerCamel, snake`)
	httpCmd.Flags().BoolVar(&doc, "doc", false, `whether generate openapi 3.0 json document or not`)
	httpCmd.Flags().StringVarP(&baseURLEnv, "env", "e", "", `base url environment variable name`)
	httpCmd.Flags().IntVarP(&routePatternStrategy, "routePattern", "r", 0, "route pattern generate strategy. 0 means splitting each methods of service interface by slash / after converting to snake case. 1 means no splitting, only lowercase. recommend default value.")
	httpCmd.Flags().BoolVar(&allowGetWithReqBody, "allowGetWithReqBody", false, "Whether allow get http request with request body.")
}
