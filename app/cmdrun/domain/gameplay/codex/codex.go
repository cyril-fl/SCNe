package codex

import (
	codextable "solopg/app/cmdrun/domain/gameplay/codex/table"
)

type Codex struct {
	Npcs      *codextable.NpcsTable
	Beastiary *codextable.BeastiaryTable
	Locations *codextable.LocationsTable
	Objects   *codextable.ObjectsTable
	Quest     *codextable.QuestTable
}

func New() *Codex {
	codex := &Codex{}
	codex.EnsureInitialized()
	return codex
}

// EnsureInitialized repairs a Codex loaded from storage where nested tables may be nil.
func (c *Codex) EnsureInitialized() {
	if c.Npcs == nil {
		c.Npcs = codextable.NewNpcs()
	}
	if c.Beastiary == nil {
		c.Beastiary = codextable.NewBeastiary()
	}
	if c.Locations == nil {
		c.Locations = codextable.NewLocations()
	}
	if c.Objects == nil {
		c.Objects = codextable.NewObjects()
	}
	if c.Quest == nil {
		c.Quest = codextable.NewQuests()
	}

	for _, table := range []codextable.Table{
		c.Npcs,
		c.Beastiary,
		c.Locations,
		c.Objects,
		c.Quest,
	} {
		table.Ensure()
	}
}
