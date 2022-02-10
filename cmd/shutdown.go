package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/cmd/internal/svc"
)

// shutdownCmd shutdowns k8s service
var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "wrap kubectl delete command to shutdown service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := svc.NewSvc("")
		s.Shutdown(k8sfile)
	},
}

func init() {
	svcCmd.AddCommand(shutdownCmd)
}
