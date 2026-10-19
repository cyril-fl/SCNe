package cmdrun

import (
	cmdrunrunsession "solopg/app/cmdrun/services/process/runsession"
	"solopg/app/shared/services/mango"
	initI19n "solopg/app/shared/services/process/initI19n"
	initlogger "solopg/app/shared/services/process/initloggerepository"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plumcleanui "github.com/cyril-fl/plum/utils/cleanui"
)

func RunCmdRun() error {
	var db *mango.Mongo
	var err error

	plumcleanui.Run()

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

	logger := initlogger.Process(db)
	translation := initI19n.Process(config.Current.I18n, "")
	session := cmdrunrunsession.Process(db)

	processes := []plumprimitive.Processable{
		logger,
		translation,
		session,
	}

	return plumprimitive.HandleProcess(processes)
}
