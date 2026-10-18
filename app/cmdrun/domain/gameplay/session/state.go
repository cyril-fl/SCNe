package session

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/cmdrun/domain/card/characters"
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/cmdrun/types/timer"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type SessionState struct {
	Player          *characters.Character
	CurrentLocation *locations.Location
	Codex           *codex.Codex
	AdventureLog    *campaign.Journal
	TimeHistory     []timer.Record
	Stopwatch       *Stopwatch
	Metadata        Metadata
}

type Metadata struct {
	Archive  plumprimitive.Datable
	Campaign plumprimitive.Datable
}

type SessionData struct {
	Campaign *campaign.Campaign
	Archives *campaign.Archives
}

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
