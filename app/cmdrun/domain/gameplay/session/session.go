package session

import (
	"scne/app/cmdrun/domain/campaign"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/card/locations"
	codextable "scne/app/cmdrun/domain/gameplay/codex/table"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/logs"
	"scne/app/shared/types/primitive"
)

// GameSession is the in-memory aggregate for an active campaign.
type GameSession struct {
	CampaignID id.ID
	State      *SessionState

	changes primitive.History[change] // this run only, not saved
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

// StartCombat sets the opponents of the combat.
func (s *GameSession) StartCombat(opponents []*characters.Character) {
	s.State.Opponents = opponents
}

// AddOpponent adds opponent to the current combat.
func (s *GameSession) AddOpponent(opponent *characters.Character) {
	s.State.Opponents = append(s.State.Opponents, opponent)
}

// EndCombat clears the opponents.
func (s *GameSession) EndCombat() {
	s.State.Opponents = nil
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

// change is a tracked player action and the player state before it.
type change struct {
	label  string
	before characters.Snapshot
}

// Track runs apply, and keeps the player state from before it so UndoLast can put it back.
// A failed apply is not kept. apply returns the chat line of the action.
func (s *GameSession) Track(apply func() (string, error)) (string, error) {
	player := s.State.Player
	if player == nil {
		return apply()
	}

	before := player.Snapshot()
	label, err := apply()
	if err == nil {
		s.changes.Push(change{label: label, before: before})
	}
	return label, err
}

// LastChange returns the chat line of the action UndoLast would undo.
func (s *GameSession) LastChange() (string, bool) {
	passed := s.changes.Passed()
	if len(passed) == 0 {
		return "", false
	}
	return passed[len(passed)-1].label, true
}

// UndoLast puts the player back as before the last tracked action and returns its chat line.
func (s *GameSession) UndoLast() (string, bool) {
	last, ok := s.changes.Undo(func(change) bool { return true })
	if !ok || s.State.Player == nil {
		return "", false
	}
	s.State.Player.Restore(last.before)
	return last.label, true
}
