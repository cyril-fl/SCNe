package form

import (
	"errors"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/types/primitive"

	tea "charm.land/bubbletea/v2"
)

// - Form - //
type Form struct {
	primitive.Fallible
	sendKey    sharedtui.EventKey
	fields     []Field
	autoSubmit bool
}

type Makefunc = func(map[string]string) *Form

func New(fields ...Field) *Form {
	return new().Add(fields...)
}

func new() *Form {
	return &Form{
		fields:     make([]Field, 0),
		autoSubmit: true,
	}
}

// Getters & Setters
func (m *Form) GetFieldValue(fieldID string) any {
	for _, f := range m.fields {
		if f.ID() == fieldID {
			return f.Value()
		}
	}

	return nil
}

func (m *Form) GetValuesAsMappedString() map[string]string {
	formValues := make(map[string]string)

	for _, f := range m.fields {
		formValues[f.ID()] = f.GetValueAsString()
	}

	return formValues
}

func (m *Form) SetSendKey(key sharedtui.EventKey) *Form {
	m.sendKey = key
	return m
}

func (m *Form) getFieldsErrors() []error {
	fieldErrors := make([]error, 0, len(m.fields))
	for _, f := range m.fields {
		if err := f.GetError(); err != nil {
			fieldErrors = append(fieldErrors, err)
		}
	}
	return fieldErrors
}

func (m *Form) getCurrentFieldIndex() int {
	for i, f := range m.fields {
		if f.IsFocused() {
			return i
		}
	}
	return -1
}

func (m *Form) setAutoSubmit(autoSubmit bool) *Form {
	m.autoSubmit = autoSubmit
	return m
}

// Methods
// - Tea Model Implementation --//
func (m *Form) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.fields))
	for _, f := range m.fields {
		cmds = append(cmds, f.Init())
	}
	return tea.Batch(cmds...)
}

func (m *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentIndex := m.getCurrentFieldIndex()
	if currentIndex == -1 {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPressMsg(msg)
	case NextField:
		if m.isLastField() && m.autoSubmit {
			return m, Send[Validate]()
		}
		return m.focusOnIndex(currentIndex + 1), nil
	case PreviousField:
		return m.focusOnIndex(currentIndex - 1), nil

	}

	return m.updateFocusedField(msg)
}

// - Fields - //
func (m *Form) Add(fields ...Field) *Form {
	m.fields = append(m.fields, fields...)
	return m
}

func (m *Form) Focus() *Form {
	return m.focusOnIndex(0)
}

func (m *Form) Submit() *Form {
	m.handleErrorsReset()
	m.Validate()

	return m
}

func (m *Form) Validate() *Form {
	m.handleErrorsReset()

	hasErrors := false
	for _, f := range m.fields {
		if err := f.Validate(); err != nil {
			f.SetError(err)
			hasErrors = true
		}
	}

	if hasErrors {
		m.SetErr(logs.Error("error.invalid:form", map[string]any{
			"Error": errors.Join(m.getFieldsErrors()...),
		}))
	}

	return m
}

func (m *Form) Reset() *Form {
	m.handleErrorsReset()

	for _, f := range m.fields {
		f.Reset()
		f.Blur()
	}

	return m.focusOnIndex(0)
}

// - Submit - //
func (m *Form) EnableAutoSubmit() *Form {
	return m.setAutoSubmit(true)
}

func (m *Form) DisableAutoSubmit() *Form {
	return m.setAutoSubmit(false)
}
