package milloracle

import (
	"testing"

	"solopg/app/cmdrun/domain/gameplay/oracle"
)

func TestRollReturnsOracleMessage(t *testing.T) {
	values := oracle.ListVisible()
	if len(values) == 0 {
		t.Skip("no visible oracles configured")
	}

	message, err := New(values[0]).Roll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == "" {
		t.Fatal("expected an oracle result message")
	}
}
