package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/enum"
)

// enumCmd updates json tag of struct fields
var enumCmd = &cobra.Command{
	Use:   "enum",
	Short: "Generate functions for constants to implement IEnum interface",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		n := enum.Generator{
			File: file,
		}
		n.Generate()
	},
}

func init() {
	rootCmd.AddCommand(enumCmd)
	enumCmd.Flags().StringVarP(&file, "file", "f", "", "absolute path of dto file")
}
