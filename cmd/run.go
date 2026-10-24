package cmd

import (
	src "solopg/app/cmdrun"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     cmd.Run.Use,
	Short:   cmd.Run.Short,
	Long:    cmd.Run.Long,
	Example: cmd.Run.Example,
	RunE:    src.RunCmdRun,
}

func init() {
	cmdRoot.AddCommand(runCmd)
}
