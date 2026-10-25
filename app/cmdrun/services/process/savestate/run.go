package savestate

import (
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/mango"
	"scne/app/shared/services/mango/repository"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type run struct {
	plumprimitive.Process

	db      *mango.Mongo
	session *session.GameSession
}

// Process returns a process saving gameSession.
func Process(db *mango.Mongo, gameSession *session.GameSession) *run {
	return &run{
		db:      db,
		session: gameSession,
	}
}

// Run saves the campaign and the archives.
func (p *run) Run() {
	p.saveCampaign()
	p.saveArchives()
}

// Getters & Setters

// GetResult is not implemented; it only logs a warning.
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
