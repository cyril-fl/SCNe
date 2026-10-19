package cmdtry

import (
	"fmt"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango"
	"solopg/app/shared/services/mango/repository"
	"solopg/app/shared/services/process/initI19n"
	"solopg/config"

	"github.com/spf13/cobra"
)

const (
	DB_ON       = false
	LOCALIZE_ON = false
	LOG_ON      = false
)

// Actuellement enregistre un log en DB
func Try(cobra *cobra.Command, args []string) error {
	var db *mango.Mongo
	var err error

	if db, err = setDb(); db != nil {
		defer db.Disconnect()
	} else if err != nil {
		return err
	}

	setI19n()
	setLogRepository(db)
	// ------------------------- //

	println("Try command executed successfully.")

	return nil
}

func setDb() (*mango.Mongo, error) {
	if !DB_ON {
		return nil, nil
	}
	return mango.Init()
}

func setLogRepository(db *mango.Mongo) {
	if !LOG_ON {
		return
	}

	if db == nil {
		panic(fmt.Errorf("database connection is nil, cannot set log repository"))
	}

	repo := repository.LogSystem()
	repo.SetDb(db)
	logs.Init(repo)
}

func setI19n() {
	if !LOCALIZE_ON {
		return
	}
	translation := initI19n.Process(config.Current.I18n, "")
	translation.Run()
}
