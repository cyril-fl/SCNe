package hintmenu

import (
	"solopg/app/cmdrun/domain/gameplay/hint"
	cmdruntui "solopg/app/cmdrun/tui"
	"solopg/app/cmdrun/tui/models"
	"solopg/app/cmdrun/tui/view/sidemenu"
	"solopg/app/cmdrun/types/direction"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// - hintMenu - //
// Menu
func NewSideMenu(size size.Size, focused bool) *HintMenu {
	items := make([]list.Item, 0)

	for _, hint := range hint.List() {
		name := i19n.Localize(hint.Name())
		items = append(items, models.NewItem(name, "", hint))
	}

	model := list.New(items, list.NewDefaultDelegate(), size.GetWidth()-4, size.GetHeight())
	models.ConfigureList(&model)
	models.SetListFocus(&model, focused)

	return &HintMenu{
		id:      "hint",
		focused: focused,
		list:    model,
	}
}

type HintMenu struct {
	id      string
	focused bool
	list    list.Model
}

// / Getters and Setters
func (m *HintMenu) ID() string {
	return m.id
}
func (m *HintMenu) IsOpen() bool {
	return false
}
func (m *HintMenu) SetOpen(open bool) {
	// No action needed as it doesn't have an open state
}

func (m *HintMenu) SetFocus(focused bool) {
	m.focused = focused
	models.SetListFocus(&m.list, focused)
}

func (m *HintMenu) GetList() *list.Model {
	return &m.list
}
func (m *HintMenu) SetList(list list.Model) {
	m.list = list
}

// Handlers
func (m *HintMenu) HandleUpdate(params sidemenu.UpdateParams) (tea.Model, tea.Cmd) {
	switch msg := params.Msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case sharedtui.CMD_SHIFT_ENTER, sharedtui.CMD_ALT_ENTER:
			return params.Model, m.handleKeyShiftEnter()
		default:
			return params.Delegate(msg)
		}
	default:
		return params.Delegate(msg)
	}
	return params.Model, nil
}

func (m *HintMenu) HandleDirectionInput(direction direction.Direction) {}

func (m *HintMenu) handleKeyShiftEnter() tea.Cmd {
	selected, ok := m.list.SelectedItem().(models.Item[hint.Hint])
	if !ok {
		return func() tea.Msg {

			return cmdruntui.SendErrorMsg(logs.Error("error.invalid", map[string]any{
				"Subject": i19n.Localize("hint"),
				"Value":   selected,
			}))
		}
	}

	spark := selected.Value()

	return func() tea.Msg {
		return Msg{
			Result: hint.Roll(spark),
		}
	}
}

// Messages
type Msg struct {
	Result []string
}
