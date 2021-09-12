package cmd

import (
	"github.com/unionj-cloud/go-doudou/svc"

	"github.com/spf13/cobra"
)

var watch bool

func newSvc() svc.Svc {
	return svc.Svc{
		Watch:      watch,
		RestartSig: make(chan int),
	}
}

// runCmd runs the service
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run go-doudou program",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		newSvc().Run()
	},
}

func init() {
	svcCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().BoolVarP(&watch, "watch", "w", false, "enable watch mode")
}
