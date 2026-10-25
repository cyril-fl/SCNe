package cmdserve

import (
	"solopg/app/cmdserve/service/process/cmdserverunsession"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango"
	"solopg/app/shared/services/mango/repository"

	"charm.land/lipgloss/v2"
	plumprimitive "github.com/cyril-fl/plum/primitive"
	"github.com/spf13/cobra"
)

func RunCmdServe(cobra *cobra.Command, args []string) error {
	db := mango.GetDb()

	log := logs.New(logs.Template{
		Type:    logs.SUCC,
		Message: "🥭 DB MongoDB connected successfully.",
	})
	lipgloss.Println(log.String())

	session := cmdserverunsession.Process(cmdserverunsession.ProcessTemplate[logs.Log]{
		GetRepository: func() repository.Watchable[logs.Log] {
			repo := repository.LogSystem()
			repo.SetDb(db)
			return repo
		},
		GetFlags: cobra.Flags,
	})

	processes := []plumprimitive.Processable{
		session,
	}

	return plumprimitive.HandleProcess(processes)
}
