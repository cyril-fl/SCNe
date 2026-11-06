package menuinventory

import (
	"os"
	"testing"

	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/card/objects"
	"scne/app/cmdrun/domain/card/objects/equipment"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/shared/components/form"
	"scne/app/shared/messages"
	"scne/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func TestMain(m *testing.M) {
	logs.Init(testLogRegistrar{})
	os.Exit(m.Run())
}

const health stats.Stat = "health"

func newPlayer(t *testing.T) *characters.Character {
	t.Helper()
	if !health.Validate() {
		t.Skip("health is not in the loaded stats config")
	}

	heal := objects.Object{Category: "potion", Effects: []stats.Effect{{Modifier: stats.Modifier{Stat: health, Value: 3}}}}
	heal.Name = "heal"
	rope := objects.Object{Category: "object"}
	rope.Name = "rope"
	helmet := equipment.Gear{Slot: "helmet"}
	helmet.Name = "helmet"
	sword := equipment.Gear{Slot: "right_hand"}
	sword.Name = "sword"

	player := &characters.Character{Stats: stats.Stats{health: 10}, Gears: []equipment.Gear{helmet}}
	player.Inventory.Add(heal, 2)
	player.Inventory.Add(rope, 1)
	player.Equipment.Equip(sword)
	return player
}

func TestMakeOptionsListsConsumablesBagAndEquipped(t *testing.T) {
	options := makeOptions(newPlayer(t))

	if len(options) != 3 {
		t.Fatalf("options = %d, want potion, bag helmet and equipped sword (no rope)", len(options))
	}
	kinds := []kind{options[0].Value().kind, options[1].Value().kind, options[2].Value().kind}
	if kinds[0] != consume || kinds[1] != equip || kinds[2] != unequip {
		t.Fatalf("option kinds = %v, want consume, equip, unequip", kinds)
	}
}

func TestApplyEachAction(t *testing.T) {
	player := newPlayer(t)
	player.Damage(health, 5)

	if _, err := apply(player, action{kind: consume, name: "heal"}); err != nil || player.Current(health) != 8 {
		t.Fatalf("consume: err = %v, health = %d; want nil, 8", err, player.Current(health))
	}
	if _, err := apply(player, action{kind: equip, index: 0}); err != nil || len(player.Gears) != 0 {
		t.Fatalf("equip: err = %v, bag = %d; want nil, 0", err, len(player.Gears))
	}
	if _, err := apply(player, action{kind: unequip, slot: "right_hand"}); err != nil || len(player.Gears) != 1 {
		t.Fatalf("unequip: err = %v, bag = %d; want nil, 1", err, len(player.Gears))
	}
	if _, err := apply(player, action{kind: unequip, slot: "boots"}); err == nil {
		t.Fatal("unequip on an empty slot, want an error")
	}
	if _, err := apply(player, action{kind: equip, index: 9}); err == nil {
		t.Fatal("equip out of the bag, want an error")
	}
}

func TestSubmitReportsResult(t *testing.T) {
	gameSession := session.NewGameSession([16]byte{}, &session.SessionState{Player: newPlayer(t)})
	f := NewForm(gameSession)
	f.Submit()

	msg := Submit(gameSession, f)()
	if _, ok := msg.(messages.Result); !ok {
		t.Fatalf("Submit message = %T, want messages.Result", msg)
	}
	if gameSession.State.Player.Inventory.Count("heal") != 1 {
		t.Fatal("Submit with the first option did not drink a potion")
	}
}

// A misclick past the last option must not apply an action: only ctrl+enter validates.
func TestFormDoesNotAutoSubmit(t *testing.T) {
	gameSession := session.NewGameSession([16]byte{}, &session.SessionState{Player: newPlayer(t)})
	f := NewForm(gameSession)

	for range 10 {
		for _, key := range []tea.KeyPressMsg{{Code: tea.KeyDown}, {Code: tea.KeyTab}, {Code: tea.KeyEnter}} {
			// Feed the produced messages back, as the bubbletea loop does.
			var msg tea.Msg = key
			for range 3 {
				_, cmd := f.Update(msg)
				if cmd == nil {
					break
				}
				if msg = cmd(); msg == nil {
					break
				}
				if _, ok := msg.(form.Validate); ok {
					t.Fatalf("%v led to form.Validate", key)
				}
			}
		}
	}
}
