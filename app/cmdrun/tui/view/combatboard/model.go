package combatboard

import (
	"solopg/app/cmdrun/domain/gameplay/session"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/types/primitive"

	tea "charm.land/bubbletea/v2"
)

// - Model - //
type model struct {
	primitive.List

	session *session.GameSession
	log     []string
	size    size.Size
}

type UiParams struct {
	Session *session.GameSession
}

func NewModel(params UiParams) *model {
	m := &model{
		session: params.Session,
	}

	m.initList()

	return m
}

// Getters & Setters

// Methods

// Update plays the selected action on ctrl+enter, leaves on esc, and moves the list otherwise.
func (m *model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	}
	return nil
}
