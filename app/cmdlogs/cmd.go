package cmdlogs

import (
	cmdlogssession "scne/app/cmdlogs/services/process/filterlogs"
	"scne/app/shared/services/mango"
	"scne/app/shared/services/mango/repository"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"github.com/spf13/cobra"
)

func RunCmdLogs(cobra *cobra.Command, args []string) error {
	db := mango.GetDb()

	repo := repository.LogSystem()
	repo.SetDb(db)

	session := cmdlogssession.Process(cmdlogssession.ProcessTemplate{
		GetList:  repo.Load,
		GetFlags: cobra.Flags,
	})

	processes := []plumprimitive.Processable{
		session,
	}

	return plumprimitive.HandleProcess(processes)
}
