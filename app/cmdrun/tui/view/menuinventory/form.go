// Package menuinventory builds the inventory form shared by the action menu and the combat.
package menuinventory

import (
	"scne/app/cmdrun/domain/card/attributes/slot"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/shared/components/form"
	"scne/app/shared/components/form/field"
	"scne/app/shared/components/models"
	"scne/app/shared/messages"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)

type kind int

const (
	consume kind = iota
	equip
	unequip
)

// action is one choice of the inventory select.
type action struct {
	kind  kind
	name  string    // consume: object name
	index int       // equip: index in the gear bag
	slot  slot.Slot // unequip
}

// NewForm returns a select listing what the player can drink, equip or take off.
func NewForm(gameSession *session.GameSession) *form.Form {
	return form.New(
		field.SelectField(field.SelectTemplate[action]{Template: field.Template[action]{ID: "action", Label: "field.inventory", Required: true}, Options: makeOptions(gameSession.State.Player)}),
	).Focus()
}

// Submit applies the chosen action to the player and reports it in the chat.
func Submit(gameSession *session.GameSession, values *form.Form) tea.Cmd {
	chosen, ok := values.GetFieldValue("action").(action)
	player := gameSession.State.Player
	if !ok || player == nil {
		return messages.Fail(logs.Error("error.unexpected:value", map[string]any{
			"Subject":  "inventory.action",
			"Expected": "action",
			"Value":    values.GetFieldValue("action"),
		}))
	}

	result, err := apply(player, chosen)
	if err != nil {
		return messages.Fail(err)
	}
	return func() tea.Msg { return messages.Result{Content: result} }
}

// apply runs chosen and returns the localized chat line.
func apply(player *characters.Character, chosen action) (string, error) {
	switch chosen.kind {
	case consume:
		if err := player.Consume(chosen.name); err != nil {
			return "", err
		}
		return i19n.Localize("inventory.msg:consumed", map[string]any{"Item": i19n.Localize(chosen.name)}), nil
	case equip:
		name := ""
		if chosen.index >= 0 && chosen.index < len(player.Gears) {
			name = player.Gears[chosen.index].GetName()
		}
		if err := player.EquipFromBag(chosen.index); err != nil {
			return "", err
		}
		return i19n.Localize("inventory.msg:equipped", map[string]any{"Item": i19n.Localize(name)}), nil
	default:
		name := equippedName(player, chosen.slot)
		if !player.Unequip(chosen.slot) {
			return "", logs.Error("error.not_found", map[string]any{"Subject": chosen.slot.String()})
		}
		return i19n.Localize("inventory.msg:unequipped", map[string]any{"Item": i19n.Localize(name)}), nil
	}
}

func makeOptions(player *characters.Character) []models.Item[action] {
	options := []models.Item[action]{}
	if player == nil {
		return options
	}

	for _, stack := range player.Inventory.Stacks() {
		if !characters.IsConsumable(stack.Object) {
			continue
		}
		label := i19n.Localize("inventory.option:consume", map[string]any{"Item": i19n.Localize(stack.GetName()), "Qty": stack.Qty})
		options = append(options, models.NewItem(label, "", action{kind: consume, name: stack.GetName()}))
	}
	for i, gear := range player.Gears {
		label := i19n.Localize("inventory.option:equip", map[string]any{"Item": i19n.Localize(gear.GetName())})
		options = append(options, models.NewItem(label, "", action{kind: equip, index: i}))
	}
	for _, gear := range player.Equipment.Gears() {
		label := i19n.Localize("inventory.option:unequip", map[string]any{"Item": i19n.Localize(gear.GetName())})
		options = append(options, models.NewItem(label, "", action{kind: unequip, slot: gear.Slot}))
	}
	return options
}

func equippedName(player *characters.Character, s slot.Slot) string {
	for _, gear := range player.Equipment.Gears() {
		if gear.Slot == s {
			return gear.GetName()
		}
	}
	return ""
}
