package savestate

import (
	"solopg/app/cmdrun/domain/gameplay/session"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango"
	"solopg/app/shared/services/mango/repository"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type run struct {
	plumprimitive.Process

	db      *mango.Mongo
	session *session.GameSession
}

func Process(db *mango.Mongo, gameSession *session.GameSession) *run {
	return &run{
		db:      db,
		session: gameSession,
	}
}

func (p *run) Run() {
	p.saveCampaign()
	p.saveArchives()
}

// Getters & Setters
func (p *run) GetResult() {
	logs.SilentWarning("error:not_implemented", map[string]any{
		"Function": "GetResult",
		"Subject":  "savestate",
	})
}

// Methods
func (p *run) saveCampaign() {
	repo := repository.Campaign().SetDb(p.db)
	data := p.session.ExportCampaign()
	p.SetErr(repo.Register(data))
}

func (p *run) saveArchives() {
	repo := repository.Archives().SetDb(p.db)
	data := p.session.ExportArchives()

	p.SetErr(repo.Register(data))
}

// Helpers
// Save returns a function that saves the game state to the database when called.
func Save(db *mango.Mongo, gameSession *session.GameSession) func() error {
	return func() error {
		process := Process(db, gameSession)
		process.Run()

		return process.GetErr()
	}
}
