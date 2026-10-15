package game

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/cmdrun/types/id"
	"solopg/app/shared/services/logs"
)

type Engine struct {
	CampaignID id.ID
	State      *State
}

func NewEngine(CampaignID id.ID, State *State) *Engine {
	return &Engine{
		CampaignID: CampaignID,
		State:      State,
	}
}

func (e *Engine) Initialize() {
	logs.SilentInfo("system.msg:init")
	e.DiscoverLocation(e.State.CurrentLocation)
}

func (e *Engine) UpdateLocation(newLocation *locations.Location) {
	e.State.CurrentLocation = newLocation
	e.DiscoverLocation(newLocation)
}

func (e *Engine) GetAdventureJournal() *campaign.Journal {
	return e.State.AdventureLog
}

func (e *Engine) DiscoverLocation(newLocation *locations.Location) {
	if newLocation == nil {
		return
	}

	key := "campaign.location:move"
	table := e.State.Codex.LocationsTable
	entry := table.FindEntryByName(newLocation.Name)

	if entry == nil {
		table.Add(newLocation)
		key = "campaign.location:discover"
	}

	logs.SilentInfo(key, map[string]any{"Location": newLocation.Name})
}

func (e *Engine) ExportArchives() *campaign.Archives {
	return &campaign.Archives{
		Datable:    e.State.Metadata.Archive.GetAge(),
		CampaignID: e.CampaignID,
		Codex:      e.State.Codex,
		Journal:    e.State.AdventureLog,
	}
}

func (e *Engine) ExportCampaign() *campaign.Campaign {
	return &campaign.Campaign{
		Datable:         e.State.Metadata.Campaign.GetAge(),
		ID:              e.CampaignID,
		Player:          e.State.Player,
		CurrentLocation: e.State.CurrentLocation,
		TimeHistory:     append(e.State.TimeHistory, e.State.Timer.GetHistory()...),
	}
}
