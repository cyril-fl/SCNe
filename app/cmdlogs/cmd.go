package cmdlogs

import (
	cmdlogssession "solopg/app/cmdlogs/services/process/filterlogs"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango"
	"solopg/app/shared/services/mango/repository"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"github.com/spf13/cobra"
)

// TEST
func RunCmdLogs(cobra *cobra.Command, args []string) error {
	db := mango.GetDb()

	session := cmdlogssession.Process(cmdlogssession.ProcessTemplate{
		GetList: func() ([]logs.Log, error) {
			repo := repository.LogSystem()
			repo.SetDb(db)
			return repo.Load(nil)
		},
		GetFlags: cobra.Flags,
	})

	processes := []plumprimitive.Processable{
		session,
	}

	return plumprimitive.HandleProcess(processes)
}
