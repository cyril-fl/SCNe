package session

import (
	"scne/app/cmdrun/domain/campaign"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/card/locations"
	"scne/app/cmdrun/domain/gameplay/codex"
	"scne/app/cmdrun/types/timer"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

// SessionState is the mutable state of a running session.
type SessionState struct {
	Player          *characters.Character
	CurrentLocation *locations.Location
	Codex           *codex.Codex
	AdventureLog    *campaign.Journal
	TimeHistory     []timer.Record
	Stopwatch       *Stopwatch
	Metadata        Metadata
	Opponents       []*characters.Character
}

// Metadata keeps the dates of the loaded documents.
type Metadata struct {
	Archive  plumprimitive.Datable
	Campaign plumprimitive.Datable
}

// SessionData holds the documents a session is loaded from.
type SessionData struct {
	Campaign *campaign.Campaign
	Archives *campaign.Archives
}

// NewSessionState builds the state from data with a new stopwatch.
func NewSessionState(data SessionData) *SessionState {
	return &SessionState{
		Player:          data.Campaign.Player,
		CurrentLocation: data.Campaign.CurrentLocation,
		Codex:           data.Archives.Codex,
		AdventureLog:    data.Archives.Journal,
		TimeHistory:     data.Campaign.TimeHistory,
		Stopwatch:       NewStopwatch(),
		Metadata: Metadata{
			Archive:  data.Archives.GetAge(),
			Campaign: data.Campaign.GetAge(),
		},
	}
}
