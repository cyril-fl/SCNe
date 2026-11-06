package characters

import (
	"testing"

	"scne/app/cmdrun/domain/card/attributes/slot"
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/objects"
	"scne/app/cmdrun/domain/card/objects/equipment"
	"scne/app/shared/services/logs"
)

const (
	health stats.Stat = "health"
	social stats.Stat = "social"
)

func requireStats(t *testing.T) {
	t.Helper()
	if !health.Validate() || !social.Validate() {
		t.Skip("health or social is not in the loaded stats config")
	}
	logs.Init(cachingTestLogRegistry{})
	t.Cleanup(func() { logs.Init(nil) })
}

func gearWith(name string, s slot.Slot, mod stats.Modifier) equipment.Gear {
	gear := equipment.Gear{Slot: s}
	gear.Name = name
	gear.Effects = []stats.Effect{{Modifier: mod}}
	return gear
}

func potion(name string, mod stats.Modifier) objects.Object {
	object := objects.Object{Category: "potion", Effects: []stats.Effect{{Modifier: mod}}}
	object.Name = name
	return object
}

func TestEffectiveStatsAddsGearWithoutChangingStats(t *testing.T) {
	requireStats(t)
	c := &Character{Stats: stats.Stats{social: 10}}
	c.Equipment.Equip(gearWith("sword", "right_hand", stats.Modifier{Stat: social, Value: 2}))

	if got := c.EffectiveStats()[social]; got != 12 {
		t.Fatalf("effective social = %d, want 12", got)
	}
	if c.Stats[social] != 10 {
		t.Fatalf("base social = %d, want 10", c.Stats[social])
	}
}

func TestDamageAndHealAreBounded(t *testing.T) {
	requireStats(t)
	c := &Character{Stats: stats.Stats{health: 10}}

	c.Damage(health, 15)
	if got := c.Current(health); got != 0 {
		t.Fatalf("health after Damage(15) = %d, want 0", got)
	}
	c.Heal(health, 4)
	if got := c.FormatStat(health); got != "4/10" {
		t.Fatalf("FormatStat(health) = %q, want 4/10", got)
	}
	c.Heal(health, 20)
	if got := c.Current(health); got != 10 {
		t.Fatalf("health after Heal(20) = %d, want 10", got)
	}
}

func TestCurrentFollowsEquipment(t *testing.T) {
	requireStats(t)
	c := &Character{Stats: stats.Stats{health: 10}}
	c.Damage(health, 3)
	c.Equipment.Equip(gearWith("amulet", "helmet", stats.Modifier{Stat: health, Value: 5}))

	if got := c.FormatStat(health); got != "12/15" {
		t.Fatalf("FormatStat(health) = %q, want 12/15", got)
	}
}

func TestEquipFromBagSwapsAndUnequip(t *testing.T) {
	requireStats(t)
	c := &Character{Gears: []equipment.Gear{gearWith("new helmet", "helmet", stats.Modifier{Stat: social, Value: 1})}}
	c.Equipment.Equip(gearWith("old helmet", "helmet", stats.Modifier{Stat: social, Value: 1}))

	if err := c.EquipFromBag(0); err != nil {
		t.Fatalf("EquipFromBag(0) error = %v", err)
	}
	if len(c.Gears) != 1 || c.Gears[0].Name != "old helmet" || c.Equipment.Gears()[0].Name != "new helmet" {
		t.Fatalf("bag = %v, equipped = %v; want old in bag, new equipped", c.Gears, c.Equipment.Gears())
	}
	if err := c.EquipFromBag(5); err == nil {
		t.Fatal("EquipFromBag(5) on a one-gear bag, want an error")
	}

	if !c.Unequip("helmet") || len(c.Gears) != 2 || len(c.Equipment.Gears()) != 0 {
		t.Fatalf("after Unequip bag = %d, equipped = %d; want 2 and 0", len(c.Gears), len(c.Equipment.Gears()))
	}
}

func TestConsume(t *testing.T) {
	requireStats(t)
	c := &Character{Stats: stats.Stats{health: 10}}
	c.Damage(health, 6)
	c.Inventory.Add(potion("heal", stats.Modifier{Stat: health, Value: 4}), 2)
	rope := objects.Object{Category: "object"}
	rope.Name = "rope"
	c.Inventory.Add(rope, 1)

	if err := c.Consume("heal"); err != nil {
		t.Fatalf("Consume(heal) error = %v", err)
	}
	if c.Current(health) != 8 || c.Inventory.Count("heal") != 1 {
		t.Fatalf("health = %d, potions = %d; want 8 and 1", c.Current(health), c.Inventory.Count("heal"))
	}
	if err := c.Consume("rope"); err == nil || c.Inventory.Count("rope") != 1 {
		t.Fatal("Consume(rope) must fail and keep the rope")
	}
	if err := c.Consume("missing"); err == nil {
		t.Fatal("Consume(missing) must fail")
	}
}

func TestPickSendsGearsToBag(t *testing.T) {
	c := &Character{}
	object := objects.Object{}
	object.Name = "thing"

	c.Pick(object, "", 2)
	c.Pick(object, "helmet", 2)

	if c.Inventory.Count("thing") != 2 || len(c.Gears) != 2 || c.Gears[0].Slot != "helmet" {
		t.Fatalf("inventory = %d, bag = %v; want 2 objects and 2 helmets", c.Inventory.Count("thing"), c.Gears)
	}
}
