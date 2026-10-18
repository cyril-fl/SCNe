package session

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/cmdrun/domain/card/locations"
	codextable "solopg/app/cmdrun/domain/gameplay/codex/table"
	"solopg/app/cmdrun/types/id"
	"solopg/app/shared/services/logs"
)

// GameSession is the in-memory aggregate for an active campaign.
type GameSession struct {
	CampaignID id.ID
	State      *SessionState
}

func NewGameSession(campaignID id.ID, state *SessionState) *GameSession {
	return &GameSession{
		CampaignID: campaignID,
		State:      state,
	}
}

func (s *GameSession) Initialize() {
	logs.SilentInfo("system.msg:init")
	s.DiscoverLocation(s.State.CurrentLocation)
}

func (s *GameSession) UpdateLocation(newLocation *locations.Location) {
	s.State.CurrentLocation = newLocation
	s.DiscoverLocation(newLocation)
}

func (s GameSession) GetAdventureJournal() *campaign.Journal {
	return s.State.AdventureLog
}

func (s *GameSession) DiscoverLocation(newLocation *locations.Location) {
	if s == nil ||
		s.State == nil ||
		s.State.Codex == nil ||
		newLocation == nil {
		return
	}

	if s.State.Codex.Locations == nil {
		s.State.Codex.EnsureInitialized()
	} else {
		s.State.Codex.Locations.Ensure()
	}

	key := "campaign.location:move"
	table := s.State.Codex.Locations

	entry := codextable.FindEntryByName(table.FormEntries(), newLocation.GetName())

	if entry == nil {
		table.Add(newLocation)
		key = "campaign.location:discover"
	}

	logs.SilentInfo(key, map[string]any{"Location": newLocation.GetName()})
}

func (s *GameSession) ExportArchives() *campaign.Archives {
	return &campaign.Archives{
		Datable:    s.State.Metadata.Archive.GetAge(),
		CampaignID: s.CampaignID,
		Codex:      s.State.Codex,
		Journal:    s.State.AdventureLog,
	}
}

func (s *GameSession) ExportCampaign() *campaign.Campaign {
	return &campaign.Campaign{
		Datable:         s.State.Metadata.Campaign.GetAge(),
		ID:              s.CampaignID,
		Player:          s.State.Player,
		CurrentLocation: s.State.CurrentLocation,
		TimeHistory:     append(s.State.TimeHistory, s.State.Stopwatch.GetHistory()...),
	}
}
