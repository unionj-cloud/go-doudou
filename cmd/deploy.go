package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"
)

// deployCmd deploy service to k8s
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "wrap command kubectl apply to deploy service to k8s",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := svc.NewSvc("")
		s.Deploy(k8sfile)
	},
}

func init() {
	svcCmd.AddCommand(deployCmd)
}
