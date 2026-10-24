package cmd

import (
	src "solopg/app/cmdtry"
	"solopg/app/shared/services/annotations"
	"solopg/config"

	"github.com/spf13/cobra"
)

var tryCmd = &cobra.Command{
	Use:         cmd.Try.Use,
	Short:       cmd.Try.Short,
	Long:        cmd.Try.Long,
	Example:     cmd.Try.Example,
	Annotations: annotations.Set(annotations.NODB, ""),
	RunE:        src.Try,
}

func init() {
	if !config.IsDev() {
		return
	}

	cmdRoot.AddCommand(tryCmd)
}
