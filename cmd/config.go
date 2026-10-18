package cmd

import (
	"solopg/config"

	plumdebug "github.com/cyril-fl/plum/utils/debug"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     cmd.Config.Use,
	Short:   cmd.Config.Short,
	Long:    cmd.Config.Long,
	Example: cmd.Config.Example,
	Run: func(cobra *cobra.Command, args []string) {
		plumdebug.Json(config.Current)
	},
}

func init() {
	cmdRoot.AddCommand(configCmd)
}
