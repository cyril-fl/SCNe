package cmd

import (
	"scne/app/cmdconfig"
	"scne/app/shared/services/annotations"

	plumcliutilsargs "github.com/cyril-fl/plum/cli/utils/args"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:         cmd.Config.Use,
	Short:       cmd.Config.Short,
	Long:        cmd.Config.Long,
	Example:     cmd.Config.Example,
	Annotations: annotations.Set(annotations.NODB, ""),
	RunE:        cmdconfig.RunECmdConfig,
}

var configEditCmd = &cobra.Command{
	Use:         cmd.ConfigEdit.Use,
	Short:       cmd.ConfigEdit.Short,
	Long:        cmd.ConfigEdit.Long,
	Example:     cmd.ConfigEdit.Example,
	Annotations: annotations.Set(annotations.NODB, ""),
	RunE:        cmdconfig.RunECmdConfigEdit,
}

func init() {
	plumcliutilsargs.Makeflags(configCmd, cmd.Config.Flags)
	configCmd.MarkFlagsMutuallyExclusive("yaml", "json")

	configCmd.AddCommand(configEditCmd)
	cmdRoot.AddCommand(configCmd)
}
