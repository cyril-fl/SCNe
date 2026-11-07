package characters

import (
	"strconv"

	"scne/app/cmdrun/domain/card/attributes/slot"
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/objects"
	"scne/app/cmdrun/domain/card/objects/equipment"
	"scne/app/shared/services/logs"
)

// shortcut: the consumable category is hard-coded, move it to category.yaml if other categories become consumable.
const consumableCategory = "potion"

// Getters & Setters

// IsConsumable reports whether object can be consumed.
func IsConsumable(object objects.Object) bool {
	return object.Category.String() == consumableCategory
}

// EffectiveStats returns the stats with the effects of the equipped gears; Stats is left untouched.
func (character *Character) EffectiveStats() stats.Stats {
	effective := make(stats.Stats, len(character.Stats))
	for stat, value := range character.Stats {
		effective[stat] = value
	}

	for _, gear := range character.Equipment.Gears() {
		for _, effect := range gear.Effects {
			effective.ApplyModifier(effect.Modifier)
		}
	}
	return effective
}

// Current returns the effective stat minus what it lost, at least 0.
func (character *Character) Current(stat stats.Stat) int {
	return max(0, character.EffectiveStats()[stat]-character.Lost[stat])
}

// FormatStat returns the effective stat, as "current/max" for health and stamina.
func (character *Character) FormatStat(stat stats.Stat) string {
	value := character.EffectiveStats()[stat]
	if !IsResource(stat) {
		return strconv.Itoa(value)
	}
	return strconv.Itoa(character.Current(stat)) + "/" + strconv.Itoa(value)
}

// IsResource reports whether stat has a current value below its maximum (health, stamina).
// shortcut: resources are recognized by name, add a "resource" flag in stats.yaml if other stats need a current value.
func IsResource(stat stats.Stat) bool {
	return stat == "health" || stat == "stamina"
}

// Methods

// IsDepleted reports whether the resource stat is at 0.
func (character *Character) IsDepleted(stat stats.Stat) bool {
	return IsResource(stat) && character.Current(stat) == 0
}

// Damage lowers stat by n, down to 0.
func (character *Character) Damage(stat stats.Stat, n int) {
	if character.Lost == nil {
		character.Lost = make(stats.Stats)
	}
	character.Lost[stat] = min(character.Lost[stat]+n, character.EffectiveStats()[stat])
}

// Heal restores stat by n, up to its effective value.
func (character *Character) Heal(stat stats.Stat, n int) {
	if character.Lost[stat] == 0 {
		return
	}
	character.Lost[stat] = max(0, character.Lost[stat]-n)
}

// EquipFromBag equips the carried gear at index i; the gear it replaces goes back to the bag.
func (character *Character) EquipFromBag(i int) error {
	if i < 0 || i >= len(character.Gears) {
		return logs.Error("error.not_found", map[string]any{
			"Subject": "gear",
		})
	}

	gear := character.Gears[i]
	character.Gears = append(character.Gears[:i], character.Gears[i+1:]...)

	if previous, had := character.Equipment.Equip(gear); had {
		character.Gears = append(character.Gears, previous)
	}
	return nil
}

// Unequip moves the gear in s to the bag; it reports false when s is empty.
func (character *Character) Unequip(s slot.Slot) bool {
	gear, had := character.Equipment.Unequip(s)
	if had {
		character.Gears = append(character.Gears, gear)
	}
	return had
}

// Pick adds qty copies of object: to the gear bag when s is a slot, to the inventory otherwise.
func (character *Character) Pick(object objects.Object, s slot.Slot, qty int) {
	if s == "" {
		character.Inventory.Add(object, qty)
		return
	}
	for range qty {
		character.Gears = append(character.Gears, equipment.Gear{Object: object, Slot: s})
	}
}

// Consume uses one object named name: positive effects heal their stat,
// negative ones damage it. Only potions can be consumed.
func (character *Character) Consume(name string) error {
	index := -1
	for i, object := range character.Inventory {
		if object.GetName() == name {
			index = i
			break
		}
	}

	if index < 0 {
		return logs.Error("error.not_found", map[string]any{
			"Subject": name,
		})
	}

	object := character.Inventory[index]
	if !IsConsumable(object) {
		return logs.Error("error.invalid", map[string]any{
			"Subject": "category",
			"Value":   object.Category,
		})
	}

	for _, effect := range object.Effects {
		if value := effect.Modifier.Value; value > 0 {
			character.Heal(effect.Modifier.Stat, value)
		} else {
			character.Damage(effect.Modifier.Stat, -value)
		}
	}

	character.Inventory.Remove(name, 1)
	return nil
}
