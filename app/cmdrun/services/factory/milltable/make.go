package milltable

import (
	codextable "solopg/app/cmdrun/domain/gameplay/codex/table"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

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
func (f mill) Make() codextable.Table {
	if *f.target == nil {
		*f.target = f.factory()
	}
	(*f.target).Ensure()
	return *f.target
}

// Helpers
func Npcs(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewNpcs()
		},
	}
}

func Beastiary(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewBeastiary()
		},
	}
}

func Locations(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewLocations()
		},
	}
}

func Objects(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewObjects()
		},
	}
}

func Quest(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewQuests()
		},
	}
}
