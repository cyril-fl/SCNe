package initloggerepository

import (
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango"
	"solopg/app/shared/services/mango/repository"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type run struct {
	plumprimitive.Process

	db *mango.Mongo
}

// Process returns a process plugging the logs into db.
func Process(db *mango.Mongo) *run {
	return &run{
		db: db,
	}
}

// Run sets the log system repository on db and initializes the logs with it.
func (p *run) Run() {
	repo := repository.LogSystem()
	repo.SetDb(p.db)

	logs.Init(repo)
}

// Getters & Setters

// Methods

// Helpers
