package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/svc"
)

// shutdownCmd shutdowns k8s service
var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "wrap kubectl delete command to shutdown service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := svc.NewSvc()
		s.Shutdown(k8sfile)
	},
}

func init() {
	svcCmd.AddCommand(shutdownCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shutdownCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
