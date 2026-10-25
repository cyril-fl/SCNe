package cmdroot

import (
	"scne/app/shared/services/annotations"
	"scne/app/shared/services/mango"
	"scne/app/shared/services/process/initI19n"
	"scne/app/shared/services/process/initloggerepository"
	"scne/app/shared/services/process/inittheme"
	"scne/config"

	"scne/app/cmdversion/services/process/generatebuildinfo"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"github.com/spf13/cobra"
)

func PersistentPreRunECmdRoot(cobra *cobra.Command, args []string) error {
	// Before the NODB check: every command renders with the theme.
	inittheme.Process(config.GetConfig().Theme).Run()

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
	translation := initI19n.Process(config.GetConfig().I18n, "")

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
