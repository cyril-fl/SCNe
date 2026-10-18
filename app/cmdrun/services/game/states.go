package game

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/cmdrun/domain/card/characters"
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/cmdrun/types/timer"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type State struct {
	Player          *characters.Character
	CurrentLocation *locations.Location
	Codex           *codex.Codex
	AdventureLog    *campaign.Journal
	TimeHistory     []timer.Record
	Timer           *Stopwatch
	Metadata        Metadata
}
type Metadata struct {
	Archive  plumprimitive.Datable
	Campaign plumprimitive.Datable
}

type CampaignData struct {
	Campaign *campaign.Campaign
	Archives *campaign.Archives
}

func NewState(data CampaignData) *State {
	return &State{
		Player:          data.Campaign.Player,
		CurrentLocation: data.Campaign.CurrentLocation,
		Codex:           data.Archives.Codex,
		AdventureLog:    data.Archives.Journal,
		TimeHistory:     data.Campaign.TimeHistory,
		Timer:           Time(),
		Metadata: Metadata{
			Archive:  data.Archives.GetAge(),
			Campaign: data.Campaign.GetAge(),
		},
	}
}
