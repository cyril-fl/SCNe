package cmd

import (
	src "solopg/app/cmdserve"

	plumcliutilsargs "github.com/cyril-fl/plum/cli/utils/args"
	"github.com/spf13/cobra"
)

var cmdServe = &cobra.Command{
	Use:     cmd.Serve.Use,
	Short:   cmd.Serve.Short,
	Long:    cmd.Serve.Long,
	Example: cmd.Serve.Example,
	RunE:    src.RunCmdServe,
}

func init() {
	plumcliutilsargs.Makeflags(cmdServe, cmd.Serve.Flags)

	cmdRoot.AddCommand(cmdServe)
}
