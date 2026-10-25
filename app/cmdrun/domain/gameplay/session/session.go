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

// NewGameSession returns a session for campaignID.
func NewGameSession(campaignID id.ID, state *SessionState) *GameSession {
	return &GameSession{
		CampaignID: campaignID,
		State:      state,
	}
}

// Initialize discovers the current location.
func (s *GameSession) Initialize() {
	logs.SilentInfo("system.msg:init")
	s.DiscoverLocation(s.State.CurrentLocation)
}

// UpdateLocation moves the player to newLocation and discovers it.
func (s *GameSession) UpdateLocation(newLocation *locations.Location) {
	s.State.CurrentLocation = newLocation
	s.DiscoverLocation(newLocation)
}

// GetAdventureJournal returns the adventure log.
func (s GameSession) GetAdventureJournal() *campaign.Journal {
	return s.State.AdventureLog
}

// DiscoverLocation adds newLocation to the codex when it is not there yet.
// It does nothing when the session, codex or location is nil.
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

// ExportArchives returns the archives to save.
func (s *GameSession) ExportArchives() *campaign.Archives {
	return &campaign.Archives{
		Datable:    s.State.Metadata.Archive.GetAge(),
		CampaignID: s.CampaignID,
		Codex:      s.State.Codex,
		Journal:    s.State.AdventureLog,
	}
}

// ExportCampaign returns the campaign to save, with the time played this session.
func (s *GameSession) ExportCampaign() *campaign.Campaign {
	return &campaign.Campaign{
		Datable:         s.State.Metadata.Campaign.GetAge(),
		ID:              s.CampaignID,
		Player:          s.State.Player,
		CurrentLocation: s.State.CurrentLocation,
		TimeHistory:     append(s.State.TimeHistory, s.State.Stopwatch.GetHistory()...),
	}
}
