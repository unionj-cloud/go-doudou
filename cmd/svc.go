package cmd

import (
	"github.com/spf13/cobra"
)

var k8sfile string
var imageRepo string

// svcCmd is the base command for generation or update
var svcCmd = &cobra.Command{
	Use:   "svc",
	Short: "generate or update service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(svcCmd)

	deployCmd.Flags().StringVarP(&k8sfile, "k8sfile", "k", "", `k8s yaml file for deploying service`)
	shutdownCmd.Flags().StringVarP(&k8sfile, "k8sfile", "k", "", `k8s yaml file for deploying service`)
	pushCmd.Flags().StringVarP(&imageRepo, "repo", "r", "", `your private docker image repository`)
}
