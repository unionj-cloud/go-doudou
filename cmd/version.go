package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v42/github"
	goversion "github.com/hashicorp/go-version"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"
	"github.com/unionj-cloud/go-doudou/v2/version"
	"os"
	"time"
)

func LatestReleaseVer() string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	release, _, err := github.NewClient(nil).Repositories.GetLatestRelease(ctx, "unionj-cloud", "go-doudou")
	if err != nil {
		panic(err)
	}
	return release.GetTagName()
}

var Prompt ISelect = &promptui.Select{
	Label:  "Do you want to upgrade?",
	Items:  []string{"Yes", "No"},
	Stdin:  os.Stdin,
	Stdout: os.Stdout,
}

var VersionSvc = svc.NewSvc
var LatestReleaseVerFunc = LatestReleaseVer

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of go-doudou",
	Long:  `You can get information about latest release version besides version number of installed go-doudou`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Installed version is %s\n", version.Release)
		latest := LatestReleaseVerFunc()
		currentVersion, _ := goversion.NewVersion(version.Release)
		latestVersion, _ := goversion.NewVersion(latest)
		if currentVersion.LessThan(latestVersion) {
			fmt.Printf("Latest release version is %s\n", latest)
			_, result, err := Prompt.Run()
			if err != nil {
				panic(err)
			}
			if result == "Yes" {
				s := VersionSvc("")
				s.Upgrade(latest)
				fmt.Println("DONE")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
