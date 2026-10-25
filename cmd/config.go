package cmd

import (
	"scne/app/cmdconfig"
	"scne/app/shared/services/annotations"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:         cmd.Config.Use,
	Short:       cmd.Config.Short,
	Long:        cmd.Config.Long,
	Example:     cmd.Config.Example,
	Annotations: annotations.Set(annotations.NODB, ""),
	Run:         cmdconfig.RunCmdConfig,
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
	configCmd.AddCommand(configEditCmd)
	cmdRoot.AddCommand(configCmd)
}
