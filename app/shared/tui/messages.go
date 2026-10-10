package sharedtui

import (
	"solopg/app/cmdrun/types/viewoptions"
	"solopg/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)

/*
BACKLOG Add Ctrl Z ctrl Y pour s'il y a des fail de et miss click
- Faire une queue Passed event
- Faire une stack future event
le tout = History {
	Passed
	Future
	Curent
}

faire une evnt undo / redo en en fonction de cmdruntui.CtrlZ / CtrlY

si on entre un nouvelle element utilser une methode dumpFuture()

gere les element en tyant les string
- Player entry
- System entry

seul le Player entry sont undoable / redoable, si on tomber sur un System entry on stop la navigation
(eviter de refaire de jeux de des, d'optimser ect)

le ctrl Z / Y ne concerne que la vue, pas les codex ect .

*/

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
