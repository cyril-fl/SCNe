package cmd

import (
	src "scne/app/cmdrun"

	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:     cmd.Play.Use,
	Aliases: []string{"run"},
	Short:   cmd.Play.Short,
	Long:    cmd.Play.Long,
	Example: cmd.Play.Example,
	RunE:    src.RunCmdRun,
}

func init() {
	cmdRoot.AddCommand(playCmd)
}
