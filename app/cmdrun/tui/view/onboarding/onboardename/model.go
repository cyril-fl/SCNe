package onboardename

import (
	"fmt"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/form/field"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	cache
	// input textinput.Model
	name string
	form *form.Form
}

type cache struct {
	size *tea.WindowSizeMsg
}

func NewModel() *model {
	return &model{
		name: "",
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
	case form.Error:
		return m.handleFormError()

	case tea.KeyMsg:
		return m.delegateInputToForm(msg)
	}
	return m, nil
}

func (m model) HandlesEscape() bool {
	return true
}

func (m *model) handleFormSubmit() (tea.Model, tea.Cmd) {
	m.form.Submit()

	if m.form.HasErr() {
		return m, form.SendErrorMsg(m.form.GetErr())
	}

	logs.SilentInfo("form.submitted", map[string]any{
		"Form": "onboardename",
	})

	if m.name = m.form.GetFieldValue("name").(string); m.name == "" {
		return m, nil
	}

	return m, sharedtui.SendResolutionMsg(sharedtui.ResolutionMsg{
		Completed: true,
		Value:     m.name,
	})
}

func (m *model) handleFormError() (tea.Model, tea.Cmd) {
	fmt.Println("Form submitted with name:", m.name)
	return m, nil
}

func makeForm() *form.Form {
	return form.
		New(
			field.TextField(field.TextTemplate[string]{ID: "name", Label: "field.name", Validator: nil, Required: true}),
		).
		SetSendKey(sharedtui.KEY_ENTER).
		Focus()
}
