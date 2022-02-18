package cmd

import (
	"github.com/unionj-cloud/go-doudou/cmd/internal/svc"

	"github.com/spf13/cobra"
)

var handler bool
var client bool
var doc bool
var jsonattrcase string
var routePatternStrategy int

// httpCmd generates scaffold code of restful service
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "generate http routes and handlers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := svc.Svc{
			Handler:              handler,
			Client:               client,
			Omitempty:            omitempty,
			Doc:                  doc,
			Jsonattrcase:         jsonattrcase,
			Env:                  baseURLEnv,
			RoutePatternStrategy: routePatternStrategy,
		}
		s.Http()
	},
}

func init() {
	svcCmd.AddCommand(httpCmd)

	httpCmd.Flags().BoolVarP(&handler, "handler", "", false, "Whether generate default handler implementation or not")
	httpCmd.Flags().BoolVarP(&client, "client", "c", false, `Whether generate default golang http client code or not`)
	httpCmd.Flags().BoolVarP(&omitempty, "omitempty", "o", false, `if true, ",omitempty" will be appended to json tag of fields in every generated anonymous struct in handlers`)
	httpCmd.Flags().StringVarP(&jsonattrcase, "case", "", "lowerCamel", `apply to json tag of fields in every generated anonymous struct in handlers. optional values: lowerCamel, snake`)
	httpCmd.Flags().BoolVarP(&doc, "doc", "", false, `whether generate openapi 3.0 json document or not`)
	httpCmd.Flags().StringVarP(&baseURLEnv, "env", "e", "", `base url environment variable name`)
	httpCmd.Flags().IntVarP(&routePatternStrategy, "routePattern", "r", 0, "route pattern generate strategy. 0 means splitting each methods of service interface by slash / after converting to snake case. 1 means no splitting, only lowercase. recommend default value.")
}
