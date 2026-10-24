package cmd

import (
	"solopg/app/cmdversion"
	"solopg/app/shared/services/annotations"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:         cmd.Version.Use,
	Short:       cmd.Version.Short,
	Long:        cmd.Version.Long,
	Example:     cmd.Version.Example,
	Annotations: annotations.Set(annotations.NODB, ""),
	Run:         cmdversion.Version,
}

func init() {
	cmdRoot.AddCommand(versionCmd)
}
