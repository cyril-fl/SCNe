package combat

import (
	"testing"
)

func TestGetShortcutsLoadsConfiguredFile(t *testing.T) {
	if got := len(GetShortcuts()); got != 10 {
		t.Fatalf("len(GetShortcuts()) = %d, want 10", got)
	}
}
