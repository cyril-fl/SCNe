package cmd

import (
	"scne/app/cmdtemplate"
	"scne/app/shared/services/annotations"

	plumcliconfig "github.com/cyril-fl/plum/cli/config"
	plumcliutilsargs "github.com/cyril-fl/plum/cli/utils/args"
	"github.com/spf13/cobra"
)

// newTemplateCmd builds a template command with one flag per template type.
func newTemplateCmd(params plumcliconfig.Cmdparams, args cobra.PositionalArgs, run func(*cobra.Command, []string) error) *cobra.Command {
	command := &cobra.Command{
		Use:         params.Use,
		Short:       params.Short,
		Long:        params.Long,
		Example:     params.Example,
		Args:        args,
		Annotations: annotations.Set(annotations.NODB, ""),
		RunE: func(command *cobra.Command, args []string) error {
			// Past argument checks, an error is about the entry, not the usage.
			command.SilenceUsage = true
			return run(command, args)
		},
	}

	plumcliutilsargs.Makeflags(command, params.Flags)
	cmdtemplate.AddKindFlags(command)
	return command
}

func init() {
	cmdRoot.AddCommand(
		newTemplateCmd(cmd.Add, cobra.NoArgs, cmdtemplate.RunECmdAdd),
		newTemplateCmd(cmd.Import, cobra.MinimumNArgs(1), cmdtemplate.RunECmdImport),
		newTemplateCmd(cmd.Export, cobra.ExactArgs(1), cmdtemplate.RunECmdExport),
		newTemplateCmd(cmd.Edit, cobra.ExactArgs(1), cmdtemplate.RunECmdEdit),
		newTemplateCmd(cmd.Delete, cobra.ExactArgs(1), cmdtemplate.RunECmdDelete),
	)
}
