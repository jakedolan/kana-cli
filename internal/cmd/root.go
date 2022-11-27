package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

var flagName string
var flagVerbose bool
var commandsRequiringSite []string

func Execute() {

	site, err := site.NewSite()
	if err != nil {
		console.Error(err, flagVerbose)
	}

	// Setup the cobra command
	cmd := &cobra.Command{
		Use:   "kana",
		Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
		Args:  cobra.NoArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			isSite, err := site.Config.ProcessNameFlag(cmd)
			if err != nil {
				console.Error(err, flagVerbose)
			}

			if !isSite && arrayContains(commandsRequiringSite, cmd.Use) {
				console.Error(fmt.Errorf("the current site you are trying to work with does not exist. Use `kana start` to initialize"), flagVerbose)
			}
		},
	}

	// Add the "name" flag to allow for sites not connected to the local directory
	cmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "Specify a name for the site, used to override using the current folder.")
	cmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Display debugging information along with detailed command output")

	// Register the subcommands
	cmd.AddCommand(
		newStartCommand(site),
		newStopCommand(site),
		newOpenCommand(site),
		newWPCommand(site),
		newDestroyCommand(site),
		newConfigCommand(site),
		newExportCommand(site),
		newVersionCommand(),
		newDbCommand(site),
	)

	// Execute anything we need to
	if err := cmd.Execute(); err != nil {
		console.Error(err, flagVerbose)
	}
}

func arrayContains(array []string, name string) bool {
	for _, value := range array {
		if value == name {
			return true
		}
	}

	return false
}
