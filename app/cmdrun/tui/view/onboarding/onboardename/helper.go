package onboardename

import (
	"scne/app/shared/components/form"
	"scne/app/shared/components/form/field"
	"scne/app/shared/messages"
	"scne/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

func makeForm() *form.Form {
	return form.
		New(
			field.TextField(field.TextTemplate[string]{ID: "name", Label: "field.name", Validator: nil, Required: true}),
		).
		SetSendKey(plumtypes.KEY_ENTER).
		Focus()
}

func (m *model) handleFormSubmit() (tea.Model, tea.Cmd) {
	m.form.Submit()

	if m.form.HasErr() {
		return m, messages.Fail(m.form.GetErr())
	}

	logs.SilentInfo("form.submitted", map[string]any{
		"Form": "onboardename",
	})

	var name string
	if name = m.form.GetFieldValue("name").(string); name == "" {
		return m, nil
	}

	return m, messages.Resolve(messages.Resolved{
		Completed: true,
		Value:     name,
	})
}

func (m *model) handleFormError(err messages.Failed) (tea.Model, tea.Cmd) {
	logs.SilentError("form.error", map[string]any{
		"Form":  "onboardename",
		"Error": err,
	})
	return m, nil
}

func (m *model) delegateInputToForm(params tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := m.form.Update(params)
	return m, cmd
}
