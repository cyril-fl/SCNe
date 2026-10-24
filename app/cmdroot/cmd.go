package cmdroot

import (
	"solopg/app/shared/services/annotations"
	"solopg/app/shared/services/mango"
	"solopg/app/shared/services/process/initI19n"
	"solopg/app/shared/services/process/initloggerepository"
	"solopg/config"

	"solopg/app/cmdversion/services/process/generatebuildinfo"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"github.com/spf13/cobra"
)

func PersistentPreRunECmdRoot(cobra *cobra.Command, args []string) error {
	if _, ok := cobra.Annotations[annotations.NODB]; ok {
		return nil
	}

	db, err := mango.Init()
	if err != nil {
		return err
	}

	return setup(db)
}

func setup(db *mango.Mongo) error {
	mango.SetDb(db)

	/* --- NOTE Everything above this line is non loggable --- */
	/*
		- STEP 1: Initialize logger repository
		- STEP 2: Initialize translation bundle
	*/

	logger := initloggerepository.Process(db)
	translation := initI19n.Process(config.Current.I18n, "")

	processes := []plumprimitive.Processable{
		logger,
		translation,
	}

	return plumprimitive.HandleProcess(processes)
}

func RunECmdRoot(cobra *cobra.Command, args []string) error {
	return cobra.Help()
}

func init() {
	generatebuildinfo.
		Process().
		Run()
}
