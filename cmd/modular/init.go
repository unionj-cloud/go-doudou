package modular

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/modular"
	"github.com/unionj-cloud/toolkit/executils"
	"github.com/unionj-cloud/toolkit/pathutils"
)

// initCmd initializes the service
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a workspace folder",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var workDir string
		if len(args) > 0 {
			workDir = args[0]
		}
		var err error
		if workDir, err = pathutils.FixPath(workDir, ""); err != nil {
			logrus.Panicln(err)
		}
		conf := modular.WorkConfig{
			WorkDir: workDir,
		}
		work := modular.NewWork(conf, executils.CmdRunner{})
		work.Init()
	},
}

func init() {
	workCmd.AddCommand(initCmd)
}
