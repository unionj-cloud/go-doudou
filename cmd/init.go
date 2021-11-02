package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/svc"
)

var modName string

// initCmd initializes the service
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a project folder",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var svcdir string
		if len(args) > 0 {
			svcdir = args[0]
		}
		var err error
		if svcdir, err = pathutils.FixPath(svcdir, ""); err != nil {
			logrus.Panicln(err)
		}
		s := svc.NewSvc(svcdir)
		s.Init()
	},
}

func init() {
	svcCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&modName, "mod", "m", "", `module name`)
}
