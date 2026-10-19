package cmdversion

import (
	"fmt"
	"solopg/app/cmdversion/services/process/generatebuildinfo"

	"github.com/spf13/cobra"
)

func Version(cobra *cobra.Command, args []string) {
	g := generatebuildinfo.Process()
	g.Run()
	i := g.GetResult()
	fmt.Println(i.String())
}
