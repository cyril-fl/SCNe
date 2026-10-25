package cmdconfig

import (
	"scne/config"

	plumdebug "github.com/cyril-fl/plum/utils/debug"

	"github.com/spf13/cobra"
)

func RunCmdConfig(cobra *cobra.Command, args []string) {
	plumdebug.Json(config.GetConfig())
}
