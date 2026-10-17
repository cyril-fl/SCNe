package game

import (
	"testing"

	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/cmdrun/domain/gameplay/codex"
	codextable "solopg/app/cmdrun/domain/gameplay/codex/table"
	"solopg/app/cmdrun/types/id"
	"solopg/app/shared/services/logs"
)

func TestNewEngineStoresCampaignAndState(t *testing.T) {
	state := &State{}
	campaignID := id.New()
	engine := NewEngine(campaignID, state)
	if engine.CampaignID != campaignID || engine.State != state {
		t.Fatal("engine did not retain its campaign ID and state")
	}
}

func TestEngineDiscoverLocationAddsLocationOnce(t *testing.T) {
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

	state := &State{Codex: &codex.Codex{
		Locations: codextable.NewLocations(nil),
	}}
	engine := NewEngine(id.New(), state)

	engine.DiscoverLocation(location)
	engine.DiscoverLocation(location)

	if entries := state.Codex.Locations.FormEntries(); len(entries) != 1 {
		t.Fatalf("location entries count = %d, want 1", len(entries))
	}
}

func TestEngineDiscoverLocationIgnoresNil(t *testing.T) {
	logs.Init(noopLogRegistry{})
	state := &State{Codex: &codex.Codex{
		Locations: codextable.NewLocations(nil),
	}}

	NewEngine(id.New(), state).DiscoverLocation(nil)

	if entries := state.Codex.Locations.FormEntries(); len(entries) != 0 {
		t.Fatalf("location entries count = %d, want 0", len(entries))
	}
}

type noopLogRegistry struct{}

func (noopLogRegistry) Register(*logs.Log) error { return nil }
