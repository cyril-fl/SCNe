package session

import (
	"errors"
	"testing"

	"scne/app/cmdrun/domain/card/attributes/rarity"
	"scne/app/cmdrun/domain/card/attributes/variety"
	"scne/app/cmdrun/domain/card/characters"
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

func TestGameSessionAddOpponentAppends(t *testing.T) {
	first, second := &characters.Character{}, &characters.Character{}
	gameSession := NewGameSession(id.New(), &SessionState{})
	gameSession.StartCombat([]*characters.Character{first})

	gameSession.AddOpponent(second)

	if len(gameSession.State.Opponents) != 2 || gameSession.State.Opponents[1] != second {
		t.Fatalf("opponents = %v, want first then second", gameSession.State.Opponents)
	}
}

func TestTrackKeepsSuccessfulActionsForUndoLast(t *testing.T) {
	player := &characters.Character{}
	gameSession := NewGameSession(id.New(), &SessionState{Player: player})

	_, _ = gameSession.Track(func() (string, error) {
		player.ToggleCondition("poisoned")
		return "poisoned", nil
	})
	_, _ = gameSession.Track(func() (string, error) { return "failed", errors.New("failed") })

	if label, ok := gameSession.LastChange(); !ok || label != "poisoned" {
		t.Fatalf("last change = %q %v, want the successful action", label, ok)
	}
	if label, ok := gameSession.UndoLast(); !ok || label != "poisoned" || len(player.Conditions) != 0 {
		t.Fatalf("undo = %q %v, conditions = %v", label, ok, player.Conditions)
	}
	if _, ok := gameSession.UndoLast(); ok {
		t.Fatal("nothing left to undo, want false")
	}
}
