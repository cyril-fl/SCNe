package gameboard

import (
	"errors"
	"testing"

	"scne/app/shared/messages"

	tea "charm.land/bubbletea/v2"
)

func TestHandleSaveInputForgetsPreviousFailure(t *testing.T) {
	m := &model{}

	m.handleSaveInput(messages.Saved{Err: errors.New("mongo down")})
	m.handleSaveInput(messages.Saved{})

	if m.HasErr() {
		t.Fatalf("successful save still reports %v", m.GetErr())
	}
}

func TestToMenuMsgRoutesArrows(t *testing.T) {
	if _, ok := toMenuMsg(tea.KeyPressMsg{Code: tea.KeyUp}); !ok {
		t.Fatal("expected plain up to reach the menu")
	}
	if _, ok := toMenuMsg(tea.KeyPressMsg{Code: tea.KeyDown, Mod: tea.ModAlt}); ok {
		t.Fatal("expected alt+down to stay in the chat input")
	}
}
