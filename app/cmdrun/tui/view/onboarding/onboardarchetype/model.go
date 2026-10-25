package onboardarchetype

import (
	"solopg/app/cmdrun/types/interfaces"
	"solopg/app/shared/components/models"
	"solopg/app/shared/messages"

	"solopg/app/shared/services/i19n"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

type step int

const (
	CHOICE_STEP step = iota
	CONFIRM_STEP
)

type model struct {
	cache
	name string
	step step
	list list.Model
}

type cache struct {
	size *tea.WindowSizeMsg
}

func NewModel[T interfaces.Archetype](name string, data []T) *model {
	items := make([]list.Item, 0, len(data))

	for _, i := range data {
		if !i.IsPlayable() {
			continue
		}

		label := plumtransform.Capitalize(i19n.Localize(i.GetName()))
		items = append(items, models.NewItem(label, "", i))
	}

	listModel := models.NewList(items, 0, 0)

	return &model{
		name: name,
		step: CHOICE_STEP,
		list: listModel,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case CHOICE_STEP:
		return m.updateList(msg)
	case CONFIRM_STEP:
		return m.updateConfirm(msg)
	default:
		return m, nil
	}
}

func (m *model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch plumtypes.Key(msg.String()) {
		case plumtypes.KEY_ENTER:
			selectedItem := m.list.SelectedItem()
			if selectedItem == nil {
				return m, nil
			}

			return m, messages.Resolve(messages.Resolved{
				Completed: true,
				Value:     selectedItem,
			})
		case plumtypes.KEY_ESC:
			// NOTE: Intercept plumtypes.KEY_ESC to prevent list.Update() from quitting.
			return m, nil

		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	// if keyMsg, ok := msg.(tea.KeyMsg); ok && plumtypes.Key(keyMsg.String()) == plumtypes.KEY_ENTER {
	// 	return m, tea.Quit
	// }
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch plumtypes.Key(msg.String()) {
		case plumtypes.KEY_ENTER:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) HandlesEscape() bool {
	return true
}
