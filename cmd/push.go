package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"
)

var imagePrefix string
var imageVer string

// pushCmd pushes image to remote docker image repository
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "wrap docker build, docker tag, docker push commands and generate or update k8s deploy yaml file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := svc.NewSvc("")
		s.Push(svc.PushConfig{
			Repo:   imageRepo,
			Prefix: imagePrefix,
			Ver:    imageVer,
		})
	},
}

func init() {
	svcCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVar(&imagePrefix, "pre", "", `image name prefix string used for building and pushing docker image`)
	pushCmd.Flags().StringVar(&imageVer, "ver", "", `docker image version`)
}
