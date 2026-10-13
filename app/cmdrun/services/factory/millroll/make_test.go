package millroll

import (
	"errors"
	"testing"

	"solopg/app/cmdrun/domain/gameplay/dice"
	"solopg/app/cmdrun/domain/gameplay/hint"
	"solopg/app/cmdrun/domain/gameplay/oracle"
	sharedtui "solopg/app/shared/tui"
)

type testRoller struct {
	content string
	err     error
}

func (r testRoller) Roll() (string, error) {
	return r.content, r.err
}

func TestExecuteRollReturnsResultMessage(t *testing.T) {
	msg := Execute(testRoller{content: "rolled"})()

	result, ok := msg.(ResultMsg)
	if !ok || result.Content != "rolled" {
		t.Fatalf("expected RollResultMsg with %q, got %#v", "rolled", msg)
	}
}

func TestExecuteRollReturnsErrorMessage(t *testing.T) {
	err := errors.New("roll failed")
	msg := Execute(testRoller{err: err})()

	result, ok := msg.(sharedtui.ErrorMsg)
	if !ok || !errors.Is(result.Err, err) {
		t.Fatalf("expected ErrorMsg with %v, got %#v", err, msg)
	}
}

func TestRollReturnsOracleMessage(t *testing.T) {
	values := oracle.ListVisible()
	if len(values) == 0 {
		t.Skip("no visible oracles configured")
	}

	message, err := Oracles(values[0]).Roll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == "" {
		t.Fatal("expected an oracle result message")
	}
}

func TestRollReturnsHintMessage(t *testing.T) {
	values := hint.List()
	if len(values) == 0 {
		t.Skip("no hints configured")
	}

	message, err := Hint(values[0]).Roll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == "" {
		t.Fatal("expected a hint result message")
	}
}

func TestRollReturnsDiceMessage(t *testing.T) {
	values := dice.List()
	if len(values) == 0 {
		t.Skip("no dice configured")
	}

	message, err := Dice(values[0]).Roll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == "" {
		t.Fatal("expected a dice result message")
	}
}
