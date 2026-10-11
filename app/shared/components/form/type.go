package form

import tea "charm.land/bubbletea/v2"

// Field represents a typed form field without exposing its concrete type.
type Field interface {
	ID() string
	Label() string
	Focus() tea.Cmd
	Blur()
	IsFocused() bool
	Value() any
	GetValueAsString() string
	Reset()
	Validate() error
	SetError(err error)
	GetError() error
	ResetError()
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
}