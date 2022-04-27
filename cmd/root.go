package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const version = "v1.0.8"

// rootCmd is the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: version,
	Use:     "go-doudou",
	Short:   "go-doudou is microservice rapid develop framework based on openapi 3.0 spec and gossip protocol",
	Long: `go-doudou works like a scaffolding tool but more than that. 
it lets api providers design their apis and help them code less. 
it generates openapi 3.0 spec json document for frontend developers or other api consumers to understand what apis there, 
consumers can import it into postman to debug and test, or upload it into some code generators to download client sdk.
it provides some useful components and middleware for constructing microservice cluster like service register and discovering, 
load balancing and so on. it just begins, more features will come out soon.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
