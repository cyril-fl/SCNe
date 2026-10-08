package models

import (
	cmdruntui "solopg/app/cmdrun/tui"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	cache
}

type cache struct {
	size *tea.WindowSizeMsg
}

func Resolve() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return cmdruntui.ResolutionMsg{
			Completed: true,
		}
	}
}

func (model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return model{}, nil
}

func (m model) View() tea.View {
	return tea.NewView("")
}
