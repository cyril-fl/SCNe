package models

import (
	sharedtui "solopg/app/shared/tui"

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
		return sharedtui.ResolutionMsg{
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
