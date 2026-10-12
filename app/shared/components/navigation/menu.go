package navigation

import (
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/types/primitive"
	"solopg/app/shared/utils/transform"

	"solopg/app/cmdrun/types/direction"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/models"
	"solopg/app/shared/services/logs"

	"charm.land/lipgloss/v2"

	sharedtui "solopg/app/shared/tui"

	"charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"
)

type Named interface {
	GetName() string
}

type Openable interface {
	IsOpen() bool
	SetOpen(bool)
}

type Directional interface {
	HandleDirectionInput(direction direction.Direction)
}

// - DiceMenu - //
type TemplateMenu[T Named] struct {
	ID               string
	Options          []T
	Focused          bool
	Footer           []string
	HandleEnterInput func(T) tea.Cmd
}

// Menu
func NewMenu[T Named](params TemplateMenu[T]) *Menu {
	items := make([]list.Item, 0, len(params.Options))

	for _, option := range params.Options {
		names := i19n.Localize(option.GetName())
		items = append(items, models.NewItem(names, "", option))
	}

	model := list.New(items, list.NewDefaultDelegate(), 0, 0)
	models.ConfigureList(&model)
	model.SetSize(0, models.ListContentHeight(&model))

	focused := params.Focused
	if !focused {
		focused = false
	}

	var handleEnterInput func() tea.Cmd
	if params.HandleEnterInput != nil {
		handleEnterInput = func() tea.Cmd {
			selected := model.SelectedItem()
			item, ok := selected.(models.Item[T])
			if !ok {
				return sharedtui.SendErrorMsg(logs.Error("error.invalid", map[string]any{
					"Subject": params.ID,
					"Value":   selected,
				}))
			}

			return params.HandleEnterInput(item.Value())
		}
	}

	return &Menu{
		Listed:            primitive.NewListed(model),
		id:                params.ID,
		focused:           focused,
		footer:            params.Footer,
		handlerEnterInput: handleEnterInput,
	}
}

type Menu struct {
	primitive.Listed

	id                string
	focused           bool
	footer            []string
	handlerEnterInput func() tea.Cmd
}

// Getters and Setters
func (m *Menu) ID() string {
	return m.id
}

func (m *Menu) SetFocus(focused bool) {
	m.focused = focused
	models.SetListFocus(m.GetList(), focused)

}
func (m *Menu) SetSize(size size.Size) {
	list := m.GetList()
	list.SetSize(size.GetWidth(), models.ListContentHeight(list))
}

// / Handlers
func (m *Menu) HandleUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch sharedtui.EventKey(msg.String()) {
		case sharedtui.CMD_SHIFT_ENTER, sharedtui.CMD_ALT_ENTER:
			if m.handlerEnterInput != nil {
				return m.handlerEnterInput()
			}
		}
	}
	return nil
}

// - View - //
func (m *Menu) GetMenuView() string {
	list := m.GetList()
	title := transform.Uppercase(i19n.Localize(m.id))
	title = lipgloss.NewStyle().Bold(true).Render(title)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		list.View(),
	)
}

func (m *Menu) GetView() string {
	return ""
}

func (m *Menu) GetFooter() []string {
	return m.footer
}
