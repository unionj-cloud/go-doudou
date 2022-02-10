package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/cmd/internal/svc"
)

var prefix string

// pushCmd pushes image to remote docker image repository
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "wrap docker build, docker tag, docker push commands and generate or update k8s deploy yaml file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := svc.NewSvc("")
		s.ImagePrefix = prefix
		s.Push(imageRepo)
	},
}

func init() {
	svcCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVar(&prefix, "pre", "", `image name prefix string used for building and pushing docker image`)
}
