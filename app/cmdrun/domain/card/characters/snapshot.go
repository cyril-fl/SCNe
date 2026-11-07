package characters

import (
	"maps"
	"slices"

	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters/wallet"
	"scne/app/cmdrun/domain/card/objects"
	"scne/app/cmdrun/domain/card/objects/equipment"
)

// Snapshot is a copy of what play changes on a character: belongings, losses and conditions.
type Snapshot struct {
	equipment  []equipment.Gear
	gears      []equipment.Gear
	inventory  objects.Inventory
	wallet     wallet.Wallet
	lost       stats.Stats
	conditions []string
}

// Snapshot copies the character state; later changes do not reach it.
func (character *Character) Snapshot() Snapshot {
	return Snapshot{
		equipment:  slices.Clone(character.Equipment.List),
		gears:      slices.Clone(character.Gears),
		inventory:  slices.Clone(character.Inventory),
		wallet:     character.Wallet,
		lost:       maps.Clone(character.Lost),
		conditions: slices.Clone(character.Conditions),
	}
}

// Restore puts back the state of s.
func (character *Character) Restore(s Snapshot) {
	character.Equipment.List = s.equipment
	character.Gears = s.gears
	character.Inventory = s.inventory
	character.Wallet = s.wallet
	character.Lost = s.lost
	character.Conditions = s.conditions
}
