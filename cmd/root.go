package cmd

import (
	"fmt"
	"os"
	"solopg/app/cmdroot"
	"solopg/app/shared/services/annotations"
	"solopg/app/shared/services/mango"
	"solopg/config"

	"github.com/spf13/cobra"
)

var c = config.Current
var cmd = c.Commands

var cmdRoot = &cobra.Command{
	Use:   c.Name,
	Short: cmd.Root.Short,
	Long:  cmd.Root.Long,
	// TODO link ca correctement
	// Version: config.Version,
	Example: cmd.Root.Example,
	PersistentPreRunE: func(cobra *cobra.Command, args []string) error {
		if _, ok := cobra.Annotations[annotations.NODB]; ok {
			return nil
		}

		db, err := mango.Init()
		if err != nil {
			return err
		}

		return cmdroot.Setup(db)

	},
	RunE: cmdroot.RunCmdRoot,
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
