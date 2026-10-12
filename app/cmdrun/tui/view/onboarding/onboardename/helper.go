package onboardename

import (
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/form/field"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"

	tea "charm.land/bubbletea/v2"
)

func makeForm() *form.Form {
	return form.
		New(
			field.TextField(field.TextTemplate[string]{ID: "name", Label: "field.name", Validator: nil, Required: true}),
		).
		SetSendKey(sharedtui.KEY_ENTER).
		Focus()
}

func (m *model) handleFormSubmit() (tea.Model, tea.Cmd) {
	m.form.Submit()

	if m.form.HasErr() {
		return m, form.SendErrorMsg(m.form.GetErr())
	}

	logs.SilentInfo("form.submitted", map[string]any{
		"Form": "onboardename",
	})

	var name string
	if name = m.form.GetFieldValue("name").(string); name == "" {
		return m, nil
	}

	return m, sharedtui.SendResolutionMsg(sharedtui.ResolutionMsg{
		Completed: true,
		Value:     name,
	})
}

func (m *model) handleFormError(err form.Error) (tea.Model, tea.Cmd) {
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
