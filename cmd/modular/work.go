package modular

import (
	"github.com/spf13/cobra"
)

// workCmd is the base command for generation or update
var workCmd = &cobra.Command{
	Use:   "work",
	Short: "Build modular application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func GetWorkCmd() *cobra.Command {
	return workCmd
}
