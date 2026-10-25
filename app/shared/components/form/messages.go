package form

import (
	tea "charm.land/bubbletea/v2"
)

// Messages

// NextField moves the focus to the next field.
type NextField struct{}

// PreviousField moves the focus to the previous field.
type PreviousField struct{}

// Validate asks the form to validate.
type Validate struct{}

// Scroll asks the parent view to scroll to the focused field.
type Scroll struct{}

// Submit signals that the form was submitted.
type Submit struct{}

// Close asks the parent view to close the form.
type Close struct{}

type msg interface {
	PreviousField |
		NextField |
		Validate |
		Scroll |
		Submit |
		Close
}

// Send returns a command emitting the form message T.
func Send[T msg]() tea.Cmd {
	return func() tea.Msg {
		return T{}
	}
}

// Open asks to open the form page with ID.
type Open struct {
	ID string
}

// OpenPage returns a command emitting Open for id.
func OpenPage(id string) tea.Cmd {
	return func() tea.Msg {
		return Open{ID: id}
	}
}
