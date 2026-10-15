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

type msg interface {
	PreviousField |
		NextField |
		Validate |
		Scroll |
		Submit
}

func Send[T msg]() tea.Cmd {
	return func() tea.Msg {
		return T{}
	}
}
