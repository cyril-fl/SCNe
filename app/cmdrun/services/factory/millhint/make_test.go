package millhint

import (
	"testing"

	"solopg/app/cmdrun/domain/gameplay/hint"
)

func TestRollReturnsHintMessage(t *testing.T) {
	values := hint.List()
	if len(values) == 0 {
		t.Skip("no hints configured")
	}

	message, err := New(values[0]).Roll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == "" {
		t.Fatal("expected a hint result message")
	}
}
