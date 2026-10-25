package generatesession

import (
	"scne/app/cmdrun/domain/campaign"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/cmdrun/services/process/resolvearchives"
	"scne/app/cmdrun/services/process/resolvecampaign"
	"scne/app/cmdrun/services/process/savestate"
	cmdruntui "scne/app/cmdrun/tui"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/mango"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type run struct {
	plumprimitive.Process
	cache

	db  *mango.Mongo
	ctx *cmdruntui.Context
}

type cache struct {
	resolvedCampaign *campaign.Campaign
	loadedArchives   *campaign.Archives
	session          *session.GameSession

	err error
}

// Process returns a process building the game session from the onboarding ctx.
func Process(db *mango.Mongo, ctx *cmdruntui.Context) *run {
	return &run{
		db:  db,
		ctx: ctx,
	}
}

// Run resolves the campaign and its archives, then starts the session.
// A new campaign is saved right away.
func (p *run) Run() {
	p.makeCampaign()
	p.makeArchives()
	p.makeAuthor()
	p.makeSession()
	p.ensureSaveState()
}

// Getters & Setters

// GetResult returns the session, nil when it could not be built.
func (p *run) GetResult() *session.GameSession {
	return p.cache.session
}

// Methods
func (p *run) makeCampaign() {
	process := resolvecampaign.Process(p.ctx)
	process.Run()

	p.cache.resolvedCampaign = process.GetResult()

	p.SetErr(process.GetErr())
}

func (p *run) makeArchives() {
	if p.HasErr() {
		return
	}

	process := resolvearchives.Process(p.db, p.cache.resolvedCampaign.ID)
	process.Run()

	p.cache.loadedArchives = process.GetResult()

	p.SetErr(process.GetErr())
}

func (p *run) makeSession() {
	if p.HasErr() {
		return
	}

	p.cache.session = session.NewGameSession(p.cache.resolvedCampaign.ID, session.NewSessionState(session.SessionData{
		Campaign: p.cache.resolvedCampaign,
		Archives: p.cache.loadedArchives,
	}))
}

func (p *run) makeAuthor() {
	logs.SetAuthor(p.cache.resolvedCampaign.ID.Base64())
}

func (p *run) ensureSaveState() {
	if p.HasErr() {
		return
	}

	p.cache.session.State.Stopwatch.Start()

	if p.ctx.SelectedSave != nil {
		return
	}

	p.cache.session.Initialize()

	process := savestate.Process(p.db, p.cache.session)
	process.Run()

	p.SetErr(process.GetErr())
}

//  Helper
