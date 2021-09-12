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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// svcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// svcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deployCmd.Flags().StringVarP(&k8sfile, "k8sfile", "k", "", `k8s yaml file for deploying service`)
	shutdownCmd.Flags().StringVarP(&k8sfile, "k8sfile", "k", "", `k8s yaml file for deploying service`)
	pushCmd.Flags().StringVarP(&imageRepo, "repo", "r", "", `your private docker image repository`)
}
