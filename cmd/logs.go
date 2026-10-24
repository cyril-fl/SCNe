package cmd

import (
	"solopg/app/cmdlogs"

	plumcliutilsargs "github.com/cyril-fl/plum/cli/utils/args"
	"github.com/spf13/cobra"
)

var cmdLogs = &cobra.Command{
	Use:     cmd.Logs.Use,
	Short:   cmd.Logs.Short,
	Long:    cmd.Logs.Long,
	Example: cmd.Logs.Example,
	RunE:    cmdlogs.RunCmdLogs,
}

func init() {
	plumcliutilsargs.Makeflags(cmdLogs, cmd.Logs.Flags)

	cmdRoot.AddCommand(cmdLogs)
}
