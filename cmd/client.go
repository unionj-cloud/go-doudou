package cmd

import (
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"

	"github.com/spf13/cobra"
)

var docfile string
var baseURLEnv string
var clientpkg string

// clientCmd generates http client code
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "generate http client from openapi 3.0 spec json file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := svc.Svc{
			DocPath:   docfile,
			Omitempty: omitempty,
			Env:       baseURLEnv,
			ClientPkg: clientpkg,
		}
		s.GenClient()
	},
}

func init() {
	httpCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&docfile, "file", "f", "", `openapi 3.0 spec json file path or download link`)
	clientCmd.Flags().StringVarP(&baseURLEnv, "env", "e", "", `base url environment variable name`)
	clientCmd.Flags().StringVarP(&clientpkg, "pkg", "p", "client", `client package name`)
	clientCmd.Flags().BoolVarP(&omitempty, "omit", "o", false, `json tag omitempty`)
}
