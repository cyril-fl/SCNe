package generatesteplist

import (
	"testing"

	"solopg/app/shared/types/step"
)

func TestGeneratorAddAndInsert(t *testing.T) {
	g := Process[string]()
	a := step.Step[string]{}
	b := step.Step[string]{}
	x := step.Step[string]{}

	g.Add(a, b).Insert(x)
	if len(g.GetSteps().Steps) != 3 {
		t.Fatalf("got %d steps, want 3", len(g.GetSteps().Steps))
	}
	if g.GetSteps().CurrentIndex != 0 || len(g.GetSteps().Steps[1:]) != 2 {
		t.Fatal("generator step list has an unexpected insertion position")
	}
}

func TestGeneratorInsertHandlesEmptyAndOutOfRangeIndex(t *testing.T) {
	g := Process[string]()
	g.Insert()
	if len(g.GetSteps().Steps) != 0 {
		t.Fatal("inserting no existing steps should be a no-op")
	}

	g.GetSteps().CurrentIndex = -1
	g.Insert(step.Step[string]{})
	if len(g.GetSteps().Steps) != 1 {
		t.Fatal("expected one inserted step")
	}
}
