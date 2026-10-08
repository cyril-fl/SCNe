package dicemenu

import (
	"solopg/app/cmdrun/domain/gameplay/dice"
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

// - DiceMenu - //
// Menu
func NewSideMenu(size size.Size, focused bool) *DiceMenu {
	options := dice.List()

	items := make([]list.Item, 0, len(options))
	for _, option := range options {
		items = append(items, models.NewItem(i19n.Localize(option.GetName()), "", option))
	}

	model := list.New(items, list.NewDefaultDelegate(), size.GetWidth()-4, size.GetHeight())
	models.ConfigureList(&model)
	models.SetListFocus(&model, focused)

	return &DiceMenu{
		id:      "dice",
		focused: focused,
		list:    model,
	}
}

type DiceMenu struct {
	id      string
	focused bool
	list    list.Model
}

// / Getters and Setters
func (m *DiceMenu) ID() string {
	return m.id
}
func (m *DiceMenu) IsOpen() bool {
	return false
}
func (m *DiceMenu) SetOpen(open bool) {
	// No action needed as it doesn't have an open state
}

func (m *DiceMenu) SetFocus(focused bool) {
	m.focused = focused
	models.SetListFocus(&m.list, focused)
}

func (m *DiceMenu) GetList() *list.Model {
	return &m.list
}
func (m *DiceMenu) SetList(list list.Model) {
	m.list = list
}

// / Handlers
func (m *DiceMenu) HandleUpdate(params sidemenu.UpdateParams) (tea.Model, tea.Cmd) {
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

func (m *DiceMenu) HandleDirectionInput(direction direction.Direction) {}

func (m *DiceMenu) handleKeyShiftEnter() tea.Cmd {
	selected, ok := m.list.SelectedItem().(models.Item[dice.Dice])
	if !ok {
		return func() tea.Msg {

			return cmdruntui.SendErrorMsg(logs.Error("error.invalid", map[string]any{
				"Subject": i19n.Localize("dice"),
				"Value":   selected,
			}))
		}
	}

	dice := selected.Value()
	return func() tea.Msg {
		return Msg{
			Dice:  dice.GetName(),
			Value: dice.Roll(),
		}
	}
}

// Messages
type Msg struct {
	Dice  string
	Value int
}
