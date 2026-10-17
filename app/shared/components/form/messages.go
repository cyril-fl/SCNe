package form

import (
	tea "charm.land/bubbletea/v2"
)

// Messages
type NextField struct{}
type PreviousField struct{}
type Validate struct{}
type Scroll struct{}
type Submit struct{}
type Close struct{}

type msg interface {
	PreviousField |
		NextField |
		Validate |
		Scroll |
		Submit |
		Close
}

func Send[T msg]() tea.Cmd {
	return func() tea.Msg {
		return T{}
	}
}

type Open struct {
	ID string
}

func OpenPage(id string) tea.Cmd {
	return func() tea.Msg {
		return Open{ID: id}
	}
}
