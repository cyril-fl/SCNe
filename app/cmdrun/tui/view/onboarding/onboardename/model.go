package onboardename

import (
	cmdruntui "solopg/app/cmdrun/tui"
	"solopg/app/shared/components/page"
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

func NewModel() model {
	input := textinput.New()
	input.Focus()

	return model{
		page: page.New(page.Template{
			Title:    "onboarding.character:create",
			Subtitle: "onboarding.character:name",
		}),

		input: input,
		name:  "",
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.input.SetWidth(msg.Width - 2)
		return m, nil
	}

	newInput, cmd := m.input.Update(msg)

	m.input = newInput
	m.name = newInput.Value()

	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == sharedtui.KEY_ENTER {
		name := strings.TrimSpace(m.name)
		if name == sharedtui.KEY_EMPTY {
			return m, nil
		}

		m.name = name

		return m, cmdruntui.SendResolutionMsg(cmdruntui.ResolutionMsg{
			Completed: true,
			Value:     m.name,
		})
	}

	return m, cmd
}
