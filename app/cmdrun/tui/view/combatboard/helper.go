package combatboard

import (
	"solopg/app/cmdrun/domain/gameplay/combat"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/models"
	"solopg/app/shared/services/i19n"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

// - Handlers - //
// Input

func (m *model) handleKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	switch plumtypes.Key(msg.String()) {
	case plumtypes.KEY_ESC:
		return m.handleKeyEsc()
	case plumtypes.CMD_CTRL_ENTER:
		return m.handleCtrlEnter()
	default:
		return m.handleModelUpdate(msg)
	}
}

func (m *model) handleKeyEsc() tea.Cmd {
	m.session.EndCombat()
	m.log = nil
	return form.Send[form.Close]()
}

// handleCtrlEnter plays one round synchronously; a refused action is shown in the log.
func (m *model) handleCtrlEnter() tea.Cmd {
	current := m.session.State.Combat
	if current == nil || current.IsOver() {
		return nil
	}

	list := m.GetList()
	selected, ok := list.SelectedItem().(models.Item[combat.Action])
	if !ok {
		return nil
	}

	events, err := current.Play(selected.Value().GetID())
	if err != nil {
		m.log = append(m.log, err.Error())
		return nil
	}

	journal := m.session.GetAdventureJournal()
	for _, event := range events {
		journal.AddEntry("System", event)
	}

	m.log = append(m.log, events...)
	return nil
}

func (m *model) handleModelUpdate(msg tea.Msg) tea.Cmd {
	list := m.GetList()
	newlist, cmd := list.Update(msg)
	m.SetList(&newlist)
	return cmd
}

// Helpers

func (m *model) initList() {
	l := m.GetList()
	*l = list.New(nil, list.NewDefaultDelegate(), 0, 0)
	models.ConfigureList(l)
}

// refresh lists the player actions of the current combat.
func (m *model) refresh() {
	l := m.GetList()
	listItems := make([]list.Item, 0)

	if current := m.session.State.Combat; current != nil {
		for _, action := range current.Player.Actions {
			listItems = append(listItems, models.NewItem(i19n.Localize(action.GetName()), "", action))
		}
	}

	l.SetItems(listItems)
}
