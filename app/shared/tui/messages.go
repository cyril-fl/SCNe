package sharedtui

import (
	"solopg/app/cmdrun/types/viewoptions"
	"solopg/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)


// Message
type ResolutionMsg struct {
	Completed bool
	Value     any
	Err       error
}

func SendResolutionMsg(params ResolutionMsg) tea.Cmd {
	return func() tea.Msg {
		return params
	}
}

// Save
type SaveMsg struct {
	Err error
}

func SendSaveMsg(save func() error) tea.Cmd {
	return func() tea.Msg {
		if save == nil {
			return SaveMsg{Err: logs.Error("error.not_found:on_save")}
		}
		return SaveMsg{Err: save()}
	}
}

// Scroll
type ScrollMsg struct {
	top    int
	height int
}

func SendSizeMsg(size *tea.WindowSizeMsg) tea.Cmd {
	return func() tea.Msg {
		return *size
	}
}

func SendScrollMsg(top, height int) tea.Cmd {
	return func() tea.Msg {
		return ScrollMsg{
			top:    top,
			height: height,
		}
	}
}
func (sm ScrollMsg) GetTop() int {
	return sm.top
}
func (sm ScrollMsg) GetHeight() int {
	return sm.height
}
func (sm ScrollMsg) GetBottom() int {
	return sm.top + sm.height
}

// Refresh
type Refresh struct {
	Option viewoptions.RefreshOption
}

func SendRefreshMsg(option viewoptions.RefreshOption) tea.Cmd {
	return func() tea.Msg {
		return Refresh{Option: option}
	}
}

// Errors
type ErrorMsg struct {
	Err error
}

func SendErrorMsg(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{Err: err}
	}
}
