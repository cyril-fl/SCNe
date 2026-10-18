package generatengine

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/cmdrun/services/game"
	"solopg/app/cmdrun/services/process/resolvearchives"
	"solopg/app/cmdrun/services/process/resolvecampaign"
	"solopg/app/cmdrun/services/process/savestate"
	cmdruntui "solopg/app/cmdrun/tui"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango"

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
	engine           *game.Engine

	err error
}

func Process(db *mango.Mongo, ctx *cmdruntui.Context) *run {
	return &run{
		db:  db,
		ctx: ctx,
	}
}

func (p *run) Run() {
	p.makeCampaign()
	p.makeArchives()
	p.makeAuthor()
	p.makeEngine()
	p.ensureSaveState()
}

// Getters & Setters
func (p *run) GetResult() *game.Engine {
	return p.cache.engine
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

func (p *run) makeEngine() {
	if p.HasErr() {
		return
	}

	p.cache.engine = game.NewEngine(p.cache.resolvedCampaign.ID, game.NewState(game.CampaignData{
		Campaign: p.cache.resolvedCampaign,
		Archives: p.cache.loadedArchives,
	}))
}

func (p *run) makeAuthor() {
	logs.SetAuthor(p.cache.resolvedCampaign.ID.String())
}

func (p *run) ensureSaveState() {
	if p.HasErr() {
		return
	}

	p.cache.engine.State.Timer.Start()

	if p.ctx.SelectedSave != nil {
		return
	}

	p.cache.engine.Initialize()

	process := savestate.Process(p.db, p.cache.engine)
	process.Run()

	p.SetErr(process.GetErr())
}

//  Helper
