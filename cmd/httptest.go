package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/cmd/internal/svc"
)

var postmanCollectionPath string
var dotenvPath string

// testCmd generates http client code
var testCmd = &cobra.Command{
	Use:   "client",
	Short: "generates integration testing code from postman collection v2.1 compatible file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := svc.Svc{
			PostmanCollectionPath: postmanCollectionPath,
			DotenvPath:            dotenvPath,
		}
		s.GenIntegrationTestingCode()
	},
}

func init() {
	httpCmd.AddCommand(testCmd)

	testCmd.Flags().StringVar(&postmanCollectionPath, "collection", "", `postman collection v2.1 compatible file disk path`)
	testCmd.Flags().StringVar(&dotenvPath, "dotenv", "", `dotenv format config file disk path only for integration testing purpose`)
	testCmd.MarkFlagRequired("collection")
	testCmd.MarkFlagRequired("dotenv")
}
