package resolvearchives

import (
	"scne/app/cmdrun/domain/campaign"
	"scne/app/cmdrun/domain/gameplay/codex"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/mango"
	"scne/app/shared/services/mango/repository"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type run struct {
	plumprimitive.Process
	cache

	db         *mango.Mongo
	campaignID id.ID
}

type cache struct {
	loadedArchives *campaign.Archives
}

// Process returns a process loading the archives of campaignID.
func Process(db *mango.Mongo, campaignID id.ID) *run {
	return &run{
		db:         db,
		campaignID: campaignID,
	}
}

// Run loads the archives, creating empty ones when none are saved.
func (p *run) Run() {
	p.loadArchives()
	p.ensureArchives()
	p.ensureArchivesInitialized()
}

// Getters & Setters

// GetResult returns the archives, nil on error.
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
	})
}

func (p *run) ensureArchivesInitialized() {
	if p.HasErr() {
		return
	}

	if p.cache.loadedArchives == nil {
		err := logs.CeaseError("error.unexpected:action", map[string]any{
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
