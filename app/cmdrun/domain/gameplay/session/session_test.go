package session

import (
	"testing"

	"scne/app/cmdrun/domain/card/attributes/rarity"
	"scne/app/cmdrun/domain/card/attributes/variety"
	"scne/app/cmdrun/domain/card/locations"
	"scne/app/cmdrun/domain/gameplay/codex"
	codextable "scne/app/cmdrun/domain/gameplay/codex/table"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/logs"
)

func TestNewGameSessionStoresCampaignAndState(t *testing.T) {
	state := &SessionState{}
	campaignID := id.New()
	gameSession := NewGameSession(campaignID, state)
	if gameSession.CampaignID != campaignID || gameSession.State != state {
		t.Fatal("game session did not retain its campaign ID and state")
	}
}

func TestGameSessionDiscoverLocationAddsLocationOnce(t *testing.T) {
	logs.Init(noopLogRegistry{})
	location, err := locations.New(locations.Template{
		Name:        "The Crossroads",
		Description: "A test location",
		Rarity:      rarity.Default(),
		Variety:     variety.MakeDefault("location_card"),
	})
	if err != nil {
		t.Fatal(err)
	}

	state := &SessionState{Codex: &codex.Codex{
		Locations: codextable.NewLocations(),
	}}
	gameSession := NewGameSession(id.New(), state)

	gameSession.DiscoverLocation(location)
	gameSession.DiscoverLocation(location)

	if entries := state.Codex.Locations.FormEntries(); len(entries) != 1 {
		t.Fatalf("location entries count = %d, want 1", len(entries))
	}
}

func TestGameSessionDiscoverLocationIgnoresNil(t *testing.T) {
	logs.Init(noopLogRegistry{})
	state := &SessionState{Codex: &codex.Codex{
		Locations: codextable.NewLocations(),
	}}

	NewGameSession(id.New(), state).DiscoverLocation(nil)

	if entries := state.Codex.Locations.FormEntries(); len(entries) != 0 {
		t.Fatalf("location entries count = %d, want 0", len(entries))
	}
}

func TestGameSessionDiscoverLocationInitializesCodexTable(t *testing.T) {
	logs.Init(noopLogRegistry{})
	state := &SessionState{Codex: &codex.Codex{
		Locations: &codextable.LocationsTable{},
	}}

	location, err := locations.New(locations.Template{
		Name:        "The Crossroads",
		Description: "A test location",
		Rarity:      rarity.Default(),
		Variety:     variety.MakeDefault("location_card"),
	})
	if err != nil {
		t.Fatal(err)
	}

	NewGameSession(id.New(), state).DiscoverLocation(location)

	if entries := state.Codex.Locations.FormEntries(); len(entries) != 1 {
		t.Fatalf("location entries count = %d, want 1", len(entries))
	}
}

type noopLogRegistry struct{}

func (noopLogRegistry) Register(*logs.Log) error { return nil }
