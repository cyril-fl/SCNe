package cmd

import (
	"fmt"
	"os"
	"solopg/app/cmdroot"
	"solopg/app/cmdversion/services/process/generatebuildinfo"
	"solopg/app/shared/services/mango"
	"solopg/config"

	"github.com/spf13/cobra"
)

var c = config.GetConfig()
var cmd = c.Commands

var cmdRoot = &cobra.Command{
	Deprecated:        "This command is deprecated and will be removed in a future version. Will be replaced by 'scne'.",
	Use:               c.Name,
	Short:             cmd.Root.Short,
	Long:              cmd.Root.Long,
	Version:           generatebuildinfo.Version,
	Example:           cmd.Root.Example,
	PersistentPreRunE: cmdroot.PersistentPreRunECmdRoot,
	RunE:              cmdroot.RunECmdRoot,
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Not PersistentPostRunE: cobra skips post-run hooks when RunE fails.
	cobra.OnFinalize(func() {
		if db := mango.GetDb(); db != nil {
			db.Disconnect()
		}
	})
}
