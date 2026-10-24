package cmd

import (
	"solopg/app/cmdconfig"
	"solopg/app/shared/services/annotations"

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

func init() {
	cmdRoot.AddCommand(configCmd)
}
