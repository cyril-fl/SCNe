package onboardename

import (
	"scne/app/shared/components/form"
	"scne/app/shared/messages"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	cache

	form *form.Form
}

type cache struct {
	size *tea.WindowSizeMsg
}

func NewModel() *model {
	return &model{
		form: makeForm(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case form.Validate:
		return m.handleFormSubmit()

	case messages.Failed:
		return m.handleFormError(msg)

	case tea.KeyMsg:
		return m.delegateInputToForm(msg)
	}
	return m, nil
}

func (m model) HandlesEscape() bool {
	return true
}
