package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/logrusorgru/aurora/v4"

	"github.com/spf13/cobra"
)

var flagConfirmDestroy bool

func newDestroyCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroys the current WordPress site. This is a permanent change.",
		Run: func(cmd *cobra.Command, args []string) {
			runDestroy(cmd, args, site)
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.Flags().BoolVar(&flagConfirmDestroy, "confirm-destroy", false, "Confirm destruction of your site (doesn't require a prompt).")

	return cmd
}

func runDestroy(cmd *cobra.Command, args []string, site *site.Site) {

	var confirmDestroy = false

	if flagConfirmDestroy {
		confirmDestroy = true
	} else {
		confirmDestroy = console.PromptConfirm(fmt.Sprintf("Are you sure you want to destroy %s? %s", aurora.Bold(aurora.Blue(site.StaticConfig.SiteName)), aurora.Bold(aurora.Yellow("This operation is destructive and cannot be undone."))), false)
	}

	if confirmDestroy {
		// Stop the WordPress site.
		err := site.StopWordPress()
		if err != nil {
			console.Error(err, flagVerbose)
		}

		// Remove the site's folder in the config directory.
		err = os.RemoveAll(site.StaticConfig.SiteDirectory)
		if err != nil {
			console.Error(err, flagVerbose)
		}

		console.Success(fmt.Sprintf("Your site, %s, has been completely destroyed.", aurora.Bold(aurora.Blue(site.StaticConfig.SiteName))))
		return
	}

	console.Error(fmt.Errorf("site destruction cancelled. No data has been lost"), flagVerbose)

}
