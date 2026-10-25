package loadsave

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/shared/components/models"
	"solopg/app/shared/messages"
	"solopg/app/shared/services/i19n"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// Methods
func (m *model) handKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch plumtypes.Key(msg.String()) {
	case plumtypes.KEY_ENTER:
		if selected, ok := m.list.SelectedItem().(models.Item[*campaign.Campaign]); ok {
			m.selected = selected.Value()
		}
		return m, func() tea.Msg {
			return messages.Resolved{
				Completed: true,
				Value:     m.selected,
			}
		}
	case plumtypes.KEY_ESC:
		// NOTE: Intercept plumtypes.KEY_ESC to prevent list.Update() from quitting.
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// Helpers
func makeItems(data []campaign.Campaign) []list.Item {
	items := []list.Item{
		models.NewItem[*campaign.Campaign](plumtransform.Capitalize(i19n.Localize("campaign:new")), "", nil),
	}

	for _, s := range data {
		newItem := models.NewItem(s.Title(), s.Description(), &s)
		items = append(items, newItem)
	}

	return items
}

func makeModel(data []campaign.Campaign) list.Model {
	items := makeItems(data)
	model := models.NewList(items, 0, 0)

	return model
}
