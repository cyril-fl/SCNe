package resolvearchives

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/cmdrun/types/id"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango"
	"solopg/app/shared/services/mango/repository"
	"solopg/app/shared/services/process"
)

type run struct {
	process.Process
	cache

	db         *mango.Mongo
	campaignID id.ID
}

type cache struct {
	loadedArchives *campaign.Archives
}

func Process(db *mango.Mongo, campaignID id.ID) *run {
	return &run{
		db:         db,
		campaignID: campaignID,
	}
}

func (p *run) Run() {
	p.loadArchives()
	p.ensureArchives()
	p.ensureArchivesInitialized()
}

// Getters & Setters
func (p *run) GetResult() *campaign.Archives {
	return p.cache.loadedArchives
}

// Methods
func (p *run) loadArchives() {
	reppo := repository.Archives()
	reppo.SetDb(p.db)

	result, err := reppo.LoadByCampaignID(p.campaignID)
	if err != nil {
		p.SetErr(err)
		return
	}

	p.cache.loadedArchives = result
}

func (p *run) ensureArchives() {
	if p.HasErr() {
		return
	}

	if p.cache.loadedArchives != nil {
		return
	}

	p.cache.loadedArchives = campaign.NewArchives(campaign.ArchivesTemplate{
		CampaignID: p.campaignID,
		Codex:      codex.New(),
		Journal:    campaign.NewJournal([]campaign.Entry{}),
		Log:        campaign.NewJournal([]campaign.Entry{}),
	})
}

func (p *run) ensureArchivesInitialized() {
	if p.HasErr() {
		return
	}

	if p.cache.loadedArchives == nil {
		err := logs.CeaseError("error:unexpected:action", map[string]any{
			"Action": "unexpected:action.load_archives",
			"Error": i19n.Localize("error.not_found", map[string]any{
				"Subject": "archives",
			}),
		})

		p.SetErr(err)
		return
	}

	p.cache.loadedArchives.EnsureInitialized()
}

// Helpers
