package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/wordpress"

	"github.com/spf13/cobra"
)

func newStopCommand(controller *docker.Controller) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops the WordPress development environment.",
		Run: func(cmd *cobra.Command, args []string) {
			runStop(cmd, args, controller)
		},
	}

	return cmd

}

func runStop(cmd *cobra.Command, args []string, controller *docker.Controller) {

	site := wordpress.NewSite(controller)

	err := site.StopWordPress()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
