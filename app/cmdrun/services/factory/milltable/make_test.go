package milltable

import (
	"testing"

	codextable "solopg/app/cmdrun/domain/gameplay/codex/table"
)

type testTable struct {
	ensured bool
}

func (t *testTable) AddFromMap(map[string]string) error {
	return nil
}

func (t *testTable) SaveFromMap(int, map[string]string) error {
	return nil
}

func (*testTable) FormEntries() []codextable.Entry {
	return nil
}

func (t *testTable) Ensure() {
	t.ensured = true
}

func TestMillTablesInitializeTargets(t *testing.T) {
	tests := []struct {
		name string
		mill func(*codextable.Table) MillTable
	}{
		{name: "npcs", mill: Npcs},
		{name: "beastiary", mill: Beastiary},
		{name: "locations", mill: Locations},
		{name: "objects", mill: Objects},
		{name: "quest", mill: Quest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := &testTable{}
			target := codextable.Table(existing)
			mill := tt.mill(&target)

			created := mill.Make()

			if created == nil {
				t.Fatal("Make returned a nil table")
			}
			if created != target {
				t.Fatal("Make returned a table different from the target")
			}
			if !existing.ensured {
				t.Fatal("Make did not ensure the target")
			}
		})
	}
}

func TestMillTableKeepsExistingTarget(t *testing.T) {
	existing := &testTable{}
	target := codextable.Table(existing)
	mill := Quest(&target)

	created := mill.Make()

	if created != target {
		t.Fatal("Make replaced an existing target")
	}
	if !existing.ensured {
		t.Fatal("Make did not ensure the existing target")
	}
}

func TestNpcsMillCreatesMissingTarget(t *testing.T) {
	var target codextable.Table

	created := Npcs(&target).Make()

	if created == nil {
		t.Fatal("Make returned a nil table")
	}
	if target == nil {
		t.Fatal("Make did not create the missing target")
	}
}
