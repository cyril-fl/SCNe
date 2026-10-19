package cmdargs

import (
	"solopg/config"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func MakeFlags(cmd *cobra.Command, args []config.Cmdarg) []*pflag.FlagSet {
	var flags []*pflag.FlagSet
	for _, arg := range args {

		switch arg.Type {
		case "int":
			v := assertIntValue(arg.Value)
			cmd.Flags().IntP(arg.Name, arg.Shorthand, v, arg.Usage)
		case "string":
			cmd.Flags().StringP(arg.Name, arg.Shorthand, arg.Value, arg.Usage)
		default:
			cmd.Flags().StringP(arg.Name, arg.Shorthand, "", arg.Usage)
		}
	}
	return flags
}

// Helpers
func assertIntValue(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return v
}
