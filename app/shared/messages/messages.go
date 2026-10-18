package messages

import (
	"solopg/app/cmdrun/types/viewoptions"
	"solopg/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)

// Message
type Delayed struct {
	Content string
}

/*
HIGH Faire des erreurs coder cfr le doc dans la branch update readme.
*/
type Failed struct {
	Err error
}

type Refreshed struct {
	Option viewoptions.RefreshOption
}

type Resolved struct {
	Completed bool
	Value     any
	Err       error
}

type Result struct {
	Content string
}

type Rolled struct{}

type Saved struct {
	Err error
}

type Scrolled struct {
	top    int
	height int
}

// Senders
func Fail(err error) tea.Cmd {
	return func() tea.Msg {
		return Failed{Err: err}
	}
}

func Resolve(params Resolved) tea.Cmd {
	return func() tea.Msg {
		return params
	}
}

func Refresh(option viewoptions.RefreshOption) tea.Cmd {
	return func() tea.Msg {
		return Refreshed{Option: option}
	}
}

func Save(save func() error) tea.Cmd {
	return func() tea.Msg {
		if save == nil {
			return Saved{Err: logs.Error("error.not_found:on_save")}
		}
		return Saved{Err: save()}
	}
}

func Resize(size *tea.WindowSizeMsg) tea.Cmd {
	return func() tea.Msg {
		return *size
	}
}

func Scroll(top, height int) tea.Cmd {
	return func() tea.Msg {
		return Scrolled{
			top:    top,
			height: height,
		}
	}
}
func (sm Scrolled) GetTop() int {
	return sm.top
}
func (sm Scrolled) GetHeight() int {
	return sm.height
}
func (sm Scrolled) GetBottom() int {
	return sm.top + sm.height
}
