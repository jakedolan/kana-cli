package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/settings"
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/logrusorgru/aurora/v4"

	"github.com/spf13/cobra"
)

var startFlags settings.StartFlags

func newStartCommand(kanaSite *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a new environment in the local folder.",
		Run: func(cmd *cobra.Command, args []string) {

			err := kanaSite.EnsureDocker()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			// Check that the site is already running and show an error if it is.
			if kanaSite.IsSiteRunning() {
				console.Error(fmt.Errorf("the site is already running. Please stop your site before running the start command"), flagVerbose)
			}

			err = kanaSite.StartSite()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			console.Success(fmt.Sprintf("Your site, %s, has has started and should be open in your default browser.", aurora.Bold(aurora.Blue(kanaSite.Settings.Name))))
		},
		Args: cobra.NoArgs,
	}

	// Add associated flags to customize the site at runtime.
	cmd.Flags().BoolVarP(&startFlags.Xdebug, "xdebug", "x", false, "Enable Xdebug when starting the container.")
	cmd.Flags().BoolVarP(&startFlags.PhpMyAdmin, "phpmyadmin", "a", false, "Enable phpMyAdmin when starting the container.")
	cmd.Flags().BoolVarP(&startFlags.IsPlugin, "plugin", "p", false, "Run the site as a plugin using the current folder as the plugin source.")
	cmd.Flags().BoolVarP(&startFlags.IsTheme, "theme", "t", false, "Run the site as a theme using the current folder as the theme source.")
	cmd.Flags().BoolVarP(&startFlags.Local, "local", "l", false, "Installs the WordPress files in your current path at ./wordpress instead of the global app path.")

	return cmd
}
