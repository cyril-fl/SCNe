package loadsave

import (
	"solopg/app/cmdrun/domain/campaign"
	cmdruntui "solopg/app/cmdrun/tui"
	"solopg/app/cmdrun/tui/models"
	"solopg/app/shared/services/i19n"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/utils/transform"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// Methods
func (m *model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.size = &msg
	return m, nil
}

func (m *model) handKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case sharedtui.KEY_ENTER:
		if selected, ok := m.list.SelectedItem().(models.Item[*campaign.Campaign]); ok {
			m.selected = selected.Value()
		}
		return m, func() tea.Msg {
			return cmdruntui.ResolutionMsg{
				Completed: true,
				Value:     m.selected,
			}
		}
		// case sharedtui.KEY_ESC:
		// TODO Implementer si besoin
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// Helpers
func makeItems(data []campaign.Campaign) []list.Item {
	items := []list.Item{
		models.NewItem[*campaign.Campaign](transform.Capitalize(i19n.Localize("campaign:new")), "", nil),
	}

	for _, s := range data {
		newItem := models.NewItem(s.Title(), s.Description(), &s)
		items = append(items, newItem)
	}

	return items
}

func makeModel(data []campaign.Campaign) list.Model {
	items := makeItems(data)
	model := list.New(items, list.NewDefaultDelegate(), 0, 0)
	models.ConfigureList(&model)

	return model
}
