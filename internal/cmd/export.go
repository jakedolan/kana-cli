package cmd

import (
	"fmt"
	"path"

	"github.com/ChrisWiegman/kana-cli/internal/app"
	"github.com/ChrisWiegman/kana-cli/internal/config"
	"github.com/ChrisWiegman/kana-cli/internal/console"

	"github.com/spf13/cobra"
)

func newExportCommand(kanaConfig *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export the current config to a .kana.json file to save with your repo.",
		Run: func(cmd *cobra.Command, args []string) {
			runExport(cmd, args, kanaConfig)
		},
		Args: cobra.ArbitraryArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.DisableFlagParsing = true

	return cmd
}

func runExport(cmd *cobra.Command, args []string, kanaConfig *config.Config) {

	site, err := app.NewSite(kanaConfig)
	if err != nil {
		console.Error(err, flagVerbose)
	}

	if !site.IsSiteRunning() {
		console.Error(fmt.Errorf("the export command only works on a running site.  Please run 'kana start' to start the site"), flagVerbose)
	}

	err = site.ExportSiteConfig()
	if err != nil {
		console.Error(err, flagVerbose)
	}

	console.Success(fmt.Sprintf("Your config has been exported to %s", path.Join(kanaConfig.Directories.Working, ".kana.json")))
}
