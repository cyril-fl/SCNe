package cmdconfig

import (
	"os"
	"scne/app/shared/services/editor"
	"scne/app/shared/services/style"
	"scne/config"

	"charm.land/lipgloss/v2"
	plumcliutilsargs "github.com/cyril-fl/plum/cli/utils/args"
	plumdebug "github.com/cyril-fl/plum/utils/debug"

	"github.com/spf13/cobra"
)

// RunECmdConfig prints the full config as YAML, or as JSON with --json; --pretty adds the theme colors.
func RunECmdConfig(cobra *cobra.Command, args []string) error {
	var asJSON, asYAML, asPretty bool
	if err := plumcliutilsargs.GetBoolFlags(cobra, map[string]*bool{"json": &asJSON, "yaml": &asYAML, "pretty": &asPretty}); err != nil {
		return err
	}

	switch {
	case asJSON:
		return printConfig(asPretty, plumdebug.Json, style.Json)
	case asYAML:
		return printConfig(asPretty, plumdebug.Yaml, style.Yaml)
	default:
		return printConfig(asPretty, plumdebug.Yaml, style.Yaml)
	}
}

// printConfig prints the config with plain, or with colored when pretty; lipgloss drops the colors when stdout is not a terminal.
func printConfig(pretty bool, plain func(any), colored func(any) (string, error)) error {
	if !pretty {
		plain(config.GetConfig())
		return nil
	}

	out, err := colored(config.GetConfig())
	if err != nil {
		return err
	}

	lipgloss.Fprintln(os.Stdout, out)
	return nil
}

// RunECmdConfigEdit opens the public config file in the user editor, then validates it.
func RunECmdConfigEdit(cobra *cobra.Command, args []string) error {
	path, err := config.PublicFile()
	if err != nil {
		return err
	}

	if err := config.EnsurePublicFile(path); err != nil {
		return err
	}

	if err := editor.Open(path); err != nil {
		return err
	}

	return config.Check()
}
