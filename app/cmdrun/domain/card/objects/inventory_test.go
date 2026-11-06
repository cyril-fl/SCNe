package objects

import (
	"testing"

	"scne/app/cmdrun/domain/card"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

func named(name string) Object {
	return Object{Card: card.Card{Descriptible: plumprimitive.Describe(plumprimitive.TemplateDescriptible{Name: name})}}
}

func TestInventoryAddCountRemove(t *testing.T) {
	var inv Inventory
	inv.Add(named("potion"), 3)
	inv.Add(named("rope"), 1)

	if got := inv.Count("potion"); got != 3 {
		t.Fatalf("Count(potion) = %d, want 3", got)
	}
	if inv.Remove("potion", 4) || inv.Count("potion") != 3 {
		t.Fatal("Remove took more than held")
	}
	if !inv.Remove("potion", 2) || inv.Count("potion") != 1 || inv.Count("rope") != 1 {
		t.Fatalf("after Remove(potion, 2) inventory = %v, want 1 potion and 1 rope", inv)
	}
}

func TestInventoryStacksKeepsFirstAppearanceOrder(t *testing.T) {
	inv := Inventory{named("rope"), named("potion"), named("rope")}

	stacks := inv.Stacks()
	if len(stacks) != 2 || stacks[0].GetName() != "rope" || stacks[0].Qty != 2 || stacks[1].Qty != 1 {
		t.Fatalf("Stacks() = %v, want rope x2 then potion x1", stacks)
	}
}
