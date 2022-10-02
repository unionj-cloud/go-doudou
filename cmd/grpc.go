package cmd

import (
	"github.com/unionj-cloud/go-doudou/cmd/internal/svc"

	"github.com/spf13/cobra"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "generate grpc service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := svc.NewSvc("")
		s.Grpc()
	},
}

func init() {
	svcCmd.AddCommand(grpcCmd)
}
