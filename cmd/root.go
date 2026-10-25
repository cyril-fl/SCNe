package cmd

import (
	"os"
	"scne/app/cmdroot"
	"scne/app/cmdversion/services/process/generatebuildinfo"
	"scne/app/shared/services/mango"
	"scne/config"

	"github.com/spf13/cobra"
)

var c = config.GetConfig()
var cmd = c.Commands

var cmdRoot = &cobra.Command{
	Use:               c.Name,
	Short:             cmd.Root.Short,
	Long:              cmd.Root.Long,
	Version:           generatebuildinfo.Version,
	Example:           cmd.Root.Example,
	PersistentPreRunE: cmdroot.PersistentPreRunECmdRoot,
	RunE:              cmdroot.RunECmdRoot,
}

func Execute() {
	// cobra already prints the error and the usage.
	if err := cmdRoot.Execute(); err != nil {
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
