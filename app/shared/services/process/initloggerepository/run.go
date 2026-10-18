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

func Process(db *mango.Mongo) *run {
	return &run{
		db: db,
	}
}

func (p *run) Run() {
	repo := repository.LogSystem()
	repo.SetDb(p.db)

	logs.Init(repo)
}

// Getters & Setters

// Methods

// Helpers
