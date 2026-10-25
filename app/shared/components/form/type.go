package form

import tea "charm.land/bubbletea/v2"

// Field represents a typed form field without exposing its concrete type.
type Field interface {
	// ID returns the field identifier, used as the value key.
	ID() string
	// Label returns the display label.
	Label() string
	// Focus gives the field keyboard focus.
	Focus() tea.Cmd
	// Blur removes the keyboard focus.
	Blur()
	// IsFocused reports whether the field has the focus.
	IsFocused() bool
	// Value returns the current typed value.
	Value() any
	// GetValueAsString returns the current value as a string.
	GetValueAsString() string
	// Reset restores the default value and clears errors.
	Reset()
	// Validate checks the current value and returns the failure, if any.
	Validate() error
	// SetError records a validation error.
	SetError(err error)
	// GetError returns the recorded errors, or nil.
	GetError() error
	// ResetError clears the recorded errors.
	ResetError()
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
}
