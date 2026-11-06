package equipment

import "testing"

func TestEquipmentEquipReplacesSlot(t *testing.T) {
	var set Equipment
	old := Gear{Slot: "helmet"}
	old.Name = "old"
	next := Gear{Slot: "helmet"}
	next.Name = "next"

	if _, had := set.Equip(old); had {
		t.Fatal("Equip on an empty slot reported a replaced gear")
	}
	previous, had := set.Equip(next)
	if !had || previous.Name != "old" {
		t.Fatalf("Equip() replaced = %q, %v; want old, true", previous.Name, had)
	}
	if gears := set.Gears(); len(gears) != 1 || gears[0].Name != "next" {
		t.Fatalf("gears = %v, want only next", gears)
	}
}

func TestEquipmentUnequipEmptySlot(t *testing.T) {
	var set Equipment
	if _, had := set.Unequip("boots"); had {
		t.Fatal("Unequip on an empty slot reported a gear")
	}
}
