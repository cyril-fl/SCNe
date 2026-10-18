package cmdlogs

import (
	cmdlogssession "solopg/app/cmdlogs/services/process/filterlogs"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango"
	"solopg/app/shared/services/mango/repository"
	"solopg/app/shared/services/process/initI19n"
	"solopg/app/shared/services/process/initloggerepository"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"github.com/spf13/pflag"
)

func RunCmdLogs(flags func() *pflag.FlagSet) error {
	var db *mango.Mongo
	var err error

	if db, err = mango.Init(); err != nil {
		return err
	}
	defer db.Disconnect()

	/* --- NOTE Everything above this line is non loggable --- */
	/*
		- STEP 1: Initialize translation bundle
		- STEP 2: Initialize logger repository
		- STEP 3: Run session
	*/

	logger := initloggerepository.Process(db)
	translation := initI19n.Process(config.Current.I18n, "")
	session := cmdlogssession.Process(cmdlogssession.ProcessTemplate{
		GetList: func() ([]logs.Log, error) {
			repo := repository.LogSystem()
			repo.SetDb(db)
			return repo.Load(nil)
		},
		GetFlags: flags,
	})

	processes := []plumprimitive.Processable{
		logger,
		translation,
		session,
	}

	return plumprimitive.HandleProcess(processes)
}
