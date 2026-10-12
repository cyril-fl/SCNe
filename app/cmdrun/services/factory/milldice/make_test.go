package milldice

import (
	"testing"

	"solopg/app/cmdrun/domain/gameplay/dice"
)

func TestRollReturnsDiceMessage(t *testing.T) {
	values := dice.List()
	if len(values) == 0 {
		t.Skip("no dice configured")
	}

	message, err := New(values[0]).Roll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == "" {
		t.Fatal("expected a dice result message")
	}
}
