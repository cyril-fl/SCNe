package characters

import (
	"slices"
	"testing"

	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters/wallet"
)

func TestRestorePutsBackTheSnapshotState(t *testing.T) {
	requireStats(t)
	c := &Character{Stats: stats.Stats{health: 10}, Wallet: wallet.NewWallet(1, 0, 0)}
	c.Equipment.Equip(gearWith("helmet", "helmet", stats.Modifier{Stat: social, Value: 1}))
	c.Equipment.Equip(gearWith("sword", "right_hand", stats.Modifier{Stat: social, Value: 1}))
	c.Inventory.Add(potion("potion", stats.Modifier{Stat: health, Value: 2}), 2)
	snapshot := c.Snapshot()

	// Unequip shifts the gears in place: the snapshot must not share them.
	c.Unequip("helmet")
	_ = c.Consume("potion")
	c.Damage(health, 4)
	_ = c.Wallet.Spend(30)
	c.ToggleCondition("poisoned")

	c.Restore(snapshot)

	if names := gearNames(c); !slices.Equal(names, []string{"helmet", "sword"}) {
		t.Fatalf("equipped = %v, want [helmet sword]", names)
	}
	if len(c.Gears) != 0 || c.Inventory.Count("potion") != 2 || c.Current(health) != 10 ||
		c.Wallet.Total() != 100 || len(c.Conditions) != 0 {
		t.Fatalf("restored character = %+v", c)
	}
}

func gearNames(c *Character) []string {
	names := []string{}
	for _, gear := range c.Equipment.Gears() {
		names = append(names, gear.Name)
	}
	return names
}
