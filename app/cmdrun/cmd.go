package cmdrun

import (
	cmdrunrunsession "scne/app/cmdrun/services/process/runsession"
	"scne/app/shared/services/mango"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plumcleanui "github.com/cyril-fl/plum/utils/cleanui"
	"github.com/spf13/cobra"
)

func RunCmdRun(cobra *cobra.Command, args []string) error {
	plumcleanui.Run()

	db := mango.GetDb()

	session := cmdrunrunsession.Process(db)

	processes := []plumprimitive.Processable{
		session,
	}

	return plumprimitive.HandleProcess(processes)
}
