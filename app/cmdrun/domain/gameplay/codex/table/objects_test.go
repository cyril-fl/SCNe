package codextable

import (
	"testing"

	"scne/app/cmdrun/domain/card/attributes/slot"
	"scne/app/cmdrun/domain/card/attributes/stats"
)

func TestObjectsSaveSlotAndEffect(t *testing.T) {
	slots, statList := slot.List(), stats.List()
	if len(slots) == 0 || len(statList) == 0 {
		t.Skip("no slot or stat in the loaded config")
	}

	table := NewObjects()
	err := table.AddFromMap(map[string]string{
		"name": "sword", "description": "sharp", "category": "weapon",
		"slot": slots[0].String(), "effect_stat": statList[0].String(), "effect_value": "2",
	})
	if err != nil {
		t.Fatalf("AddFromMap error = %v", err)
	}

	object, s, found := table.Find(table.Entries[0].ID)
	if !found || s != slots[0] || len(object.Effects) != 1 || object.Effects[0].Modifier.Value != 2 {
		t.Fatalf("Find() = %+v, %q, %v; want the slot and a +2 effect", object, s, found)
	}

	values := table.FormEntries()[0].Values
	if values["slot"] != slots[0].String() || values["effect_value"] != "2" {
		t.Fatalf("FormEntries values = %v, want slot and effect for editing", values)
	}
}

func TestObjectsRejectUnknownSlot(t *testing.T) {
	table := NewObjects()
	if err := table.AddFromMap(map[string]string{"name": "n", "description": "d", "slot": "tail"}); err == nil {
		t.Fatal("AddFromMap with an unknown slot, want an error")
	}
}
