package oraclemenu

import (
	"fmt"
	"solopg/app/cmdrun/domain/gameplay/oracle"
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

// - OracleMenu - //
// Menu
func NewSideMenu(size size.Size, focused bool) *OracleMenu {
	items := make([]list.Item, 0)

	for _, rules := range oracle.List() {
		if rules.IsVisible() {
			key := fmt.Sprintf("oracle.%s", rules.Name())
			name := i19n.Localize(key)
			items = append(items, models.NewItem(name, "", rules))
		}
	}

	model := list.New(items, list.NewDefaultDelegate(), size.GetWidth()-4, size.GetHeight())
	models.ConfigureList(&model)
	models.SetListFocus(&model, focused)

	return &OracleMenu{
		id:      "oracles",
		focused: focused,
		list:    model,
	}
}

type OracleMenu struct {
	id      string
	focused bool
	list    list.Model
}

// / Getters and Setters
func (m *OracleMenu) ID() string {
	return m.id
}
func (m *OracleMenu) IsOpen() bool {
	return false
}
func (m *OracleMenu) SetOpen(open bool) {
	// No action needed as it doesn't have an open state
}

func (m *OracleMenu) SetFocus(focused bool) {
	m.focused = focused
	models.SetListFocus(&m.list, focused)
}

func (m *OracleMenu) GetList() *list.Model {
	return &m.list
}
func (m *OracleMenu) SetList(list list.Model) {
	m.list = list
}

// Handlers
func (m *OracleMenu) HandleUpdate(params sidemenu.UpdateParams) (tea.Model, tea.Cmd) {
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

func (m *OracleMenu) HandleDirectionInput(direction direction.Direction) {}

func (m *OracleMenu) handleKeyShiftEnter() tea.Cmd {
	selected, ok := m.list.SelectedItem().(models.Item[oracle.Oracle])

	if !ok {
		return func() tea.Msg {

			return cmdruntui.SendErrorMsg(logs.Error("error.invalid", map[string]any{
				"Subject": "oracle",
				"Value":   selected,
			}))
		}
	}

	rules := selected.Value()

	rollResult, err := oracle.Roll[any](rules)
	if err != nil {
		return func() tea.Msg {

			return cmdruntui.SendErrorMsg(logs.Error("error.unexpected:action:oracle", map[string]any{
				"Action": "unexpected:action.roll",
				"Error":  err,
			}))
		}
	}

	return func() tea.Msg {
		return Msg{
			Result: rollResult,
		}
	}
}

// Messages
type Msg struct {
	Result *oracle.Result[any]
}
