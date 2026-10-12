package onboardarchetype

import (
	"solopg/app/cmdrun/domain/card/attributes/archetypes"
	"solopg/app/shared/components/models"

	"solopg/app/shared/services/i19n"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/utils/transform"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
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

func NewModel[T archetypes.Archetype](name string, data []T) *model {
	items := make([]list.Item, 0, len(data))

	for _, i := range data {
		if !i.IsPlayable() {
			continue
		}

		label := transform.Capitalize(i19n.Localize(i.GetName()))
		items = append(items, models.NewItem(label, "", i))
	}

	listModel := list.New(items, list.NewDefaultDelegate(), 0, 0)
	models.ConfigureList(&listModel)

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
		switch sharedtui.EventKey(msg.String()) {
		case sharedtui.KEY_ENTER:
			selectedItem := m.list.SelectedItem()
			if selectedItem == nil {
				return m, nil
			}

			return m, sharedtui.SendResolutionMsg(sharedtui.ResolutionMsg{
				Completed: true,
				Value:     selectedItem,
			})
		case sharedtui.KEY_ESC:
			// NOTE: Intercept sharedtui.KEY_ESC to prevent list.Update() from quitting.
			return m, nil

		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	// if keyMsg, ok := msg.(tea.KeyMsg); ok && sharedtui.EventKey(keyMsg.String()) == sharedtui.KEY_ENTER {
	// 	return m, tea.Quit
	// }
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch sharedtui.EventKey(msg.String()) {
		case sharedtui.KEY_ENTER:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) HandlesEscape() bool {
	return true
}
