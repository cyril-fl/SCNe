package gameboard

import (
	"errors"
	"testing"

	"scne/app/shared/messages"
)

func TestHandleSaveInputForgetsPreviousFailure(t *testing.T) {
	m := &model{}

	m.handleSaveInput(messages.Saved{Err: errors.New("mongo down")})
	m.handleSaveInput(messages.Saved{})

	if m.HasErr() {
		t.Fatalf("successful save still reports %v", m.GetErr())
	}
}
