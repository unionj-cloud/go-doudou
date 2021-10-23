package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/svc"
)

// deployCmd deploy service to k8s
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "wrap command kubectl apply to deploy service to k8s",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := svc.NewSvc()
		s.Deploy(k8sfile)
	},
}

func init() {
	svcCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
