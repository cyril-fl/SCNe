package milltable

import (
	codextable "scne/app/cmdrun/domain/gameplay/codex/table"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

// MillTable builds a Codex table on demand.
type MillTable interface {
	Make() codextable.Table
}

type mill struct {
	plumprimitive.Mill

	target  *codextable.Table
	factory func() codextable.Table
}

// Getters & Setters

// Methods

// Make creates the target table when nil, ensures it and returns it.
func (f mill) Make() codextable.Table {
	if *f.target == nil {
		*f.target = f.factory()
	}
	(*f.target).Ensure()
	return *f.target
}

// Helpers

// Npcs returns a MillTable for the NPC table stored in target.
func Npcs(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewNpcs()
		},
	}
}

// Beastiary returns a MillTable for the beastiary table stored in target.
func Beastiary(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewBeastiary()
		},
	}
}

// Locations returns a MillTable for the locations table stored in target.
func Locations(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewLocations()
		},
	}
}

// Objects returns a MillTable for the objects table stored in target.
func Objects(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewObjects()
		},
	}
}

// Quest returns a MillTable for the quest table stored in target.
func Quest(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewQuests()
		},
	}
}
