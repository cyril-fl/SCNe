package loadsave

import (
	"scne/app/cmdrun/domain/campaign"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	list     list.Model
	selected *campaign.Campaign
}

func NewModel(saves []campaign.Campaign) *model {
	return &model{
		list: makeModel(saves),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handKeyMsg(msg)
	default:
		return m, nil
	}
}

// HandlesEscape lets the global router forward Escape while a Codex page is open.
func (m model) HandlesEscape() bool {
	return true
}
