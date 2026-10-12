package gameboard

import (
	"errors"
	"testing"

	"solopg/app/cmdrun/services/factory/roll"
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
	msg := roll.Execute(testRoller{content: "rolled"})()

	result, ok := msg.(roll.RollResultMsg)
	if !ok || result.Content != "rolled" {
		t.Fatalf("expected RollResultMsg with %q, got %#v", "rolled", msg)
	}
}

func TestExecuteRollReturnsErrorMessage(t *testing.T) {
	err := errors.New("roll failed")
	msg := roll.Execute(testRoller{err: err})()

	result, ok := msg.(sharedtui.ErrorMsg)
	if !ok || !errors.Is(result.Err, err) {
		t.Fatalf("expected ErrorMsg with %v, got %#v", err, msg)
	}
}
