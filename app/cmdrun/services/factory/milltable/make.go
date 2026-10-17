package milltable

import (
	codextable "solopg/app/cmdrun/domain/gameplay/codex/table"
	"solopg/app/shared/services/factory"
)

// TEST
type MillTable interface {
	Make() codextable.Table
}

type mill struct {
	factory.Mill

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
			return codextable.NewNpcsTable(nil)
		},
	}
}

func Beastiary(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewBeastiaryTable(nil)
		},
	}
}

func Locations(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewLocationsTable(nil)
		},
	}
}

func Objects(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewObjectsTable(nil)
		},
	}
}

func Quest(target *codextable.Table) MillTable {
	return mill{
		target: target,
		factory: func() codextable.Table {
			return codextable.NewQuestTable(nil)
		},
	}
}
