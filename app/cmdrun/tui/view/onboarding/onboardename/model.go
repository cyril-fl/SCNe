package onboardename

import (
	sharedtui "solopg/app/shared/tui"
	interfass "solopg/app/shared/types/interface"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	cache
	page interfass.Page

	input textinput.Model
	name  string
}

type cache struct {
	size *tea.WindowSizeMsg
}

func NewModel() *model {
	input := textinput.New()
	input.Focus()

	return &model{
		input: input,
		name:  "",
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case sharedtui.KEY_ENTER:
			name := strings.TrimSpace(m.name)
			if name == sharedtui.KEY_EMPTY {
				return m, nil
			}
			m.name = name

			return m, sharedtui.SendResolutionMsg(sharedtui.ResolutionMsg{
				Completed: true,
				Value:     m.name,
			})
		default:
			return m.handleDefaultInput(msg)
		}
	}
	return m, nil
}

func (m model) HandlesEscape() bool {
	return true
}
