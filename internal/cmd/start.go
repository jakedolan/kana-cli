package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/internal/traefik"

	"github.com/spf13/cobra"
)

var flagXdebug bool
var flagLocal bool
var flagIsTheme bool
var flagIsPlugin bool

func newStartCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a new environment in the local folder.",
		Run: func(cmd *cobra.Command, args []string) {
			runStart(cmd, args, site)
		},
		Args: cobra.NoArgs,
	}

	// Add associated flags to customize the site at runtime.
	cmd.Flags().BoolVarP(&flagXdebug, "xdebug", "x", false, "Enable Xdebug when starting the container.")
	cmd.Flags().BoolVarP(&flagIsPlugin, "plugin", "p", false, "Run the site as a plugin using the current folder as the plugin source.")
	cmd.Flags().BoolVarP(&flagIsTheme, "theme", "t", false, "Run the site as a theme using the current folder as the theme source.")
	cmd.Flags().BoolVarP(&flagLocal, "local", "l", false, "Installs the WordPress files in your current path at ./wordpress instead of the global app path.")

	return cmd
}

func runStart(cmd *cobra.Command, args []string, kanaSite *site.Site) {

	// A site shouldn't be both a plugin and a theme so this reports an error if that is the case.
	if flagIsPlugin && flagIsTheme {
		console.Error(fmt.Errorf("you have set both the plugin and theme flags. Please choose only one option"), flagVerbose)
	}

	// Check that the site is already running and show an error if it is.
	if kanaSite.IsSiteRunning() {
		console.Error(fmt.Errorf("the site is already running. Please stop your site before running the start command"), flagVerbose)
	}

	// Process any overrides set with flags on the start command
	startFlags := site.SiteFlags{
		Xdebug:   flagXdebug,
		IsTheme:  flagIsTheme,
		IsPlugin: flagIsPlugin,
		Local:    flagLocal,
	}

	kanaSite.ProcessSiteFlags(cmd, startFlags)

	// Let's start everything up
	fmt.Printf("Starting development site: %s\n", kanaSite.GetURL(false))

	// Start Traefik if we need it
	traefikClient, err := traefik.NewTraefik(kanaSite.StaticConfig)
	if err != nil {
		console.Error(err, flagVerbose)
	}

	err = traefikClient.StartTraefik()
	if err != nil {
		console.Error(err, flagVerbose)
	}

	// Start WordPress
	err = kanaSite.StartWordPress()
	if err != nil {
		console.Error(err, flagVerbose)
	}

	// Make sure the WordPress site is running
	_, err = kanaSite.VerifySite()
	if err != nil {
		console.Error(err, flagVerbose)
	}

	// Setup WordPress
	err = kanaSite.InstallWordPress()
	if err != nil {
		console.Error(err, flagVerbose)
	}

	// Install Xdebug if we need to
	_, err = kanaSite.InstallXdebug()
	if err != nil {
		console.Error(err, flagVerbose)
	}

	// Install any configuration plugins if needed
	err = kanaSite.InstallDefaultPlugins()
	if err != nil {
		console.Error(err, flagVerbose)
	}

	// Open the site in the user's browser
	err = kanaSite.OpenSite()
	if err != nil {
		console.Error(err, flagVerbose)
	}

	console.Success("Your site has started and should be open in your default browser.")
}
