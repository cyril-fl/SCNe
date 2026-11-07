package characters

import (
	"slices"
	"testing"
)

func TestToggleConditionAddsThenRemoves(t *testing.T) {
	c := &Character{}

	if !c.ToggleCondition(" poisoned ") || !slices.Equal(c.Conditions, []string{"poisoned"}) {
		t.Fatalf("conditions = %v, want [poisoned]", c.Conditions)
	}
	if c.ToggleCondition("POISONED") || len(c.Conditions) != 0 {
		t.Fatalf("toggling again ignoring case must remove it, conditions = %v", c.Conditions)
	}
	if c.ToggleCondition("  ") || len(c.Conditions) != 0 {
		t.Fatalf("a blank condition must change nothing, conditions = %v", c.Conditions)
	}
}
