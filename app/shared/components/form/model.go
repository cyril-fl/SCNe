package form

import (
	"errors"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"

	tea "charm.land/bubbletea/v2"
)

// - Form - //
type Form struct {
	fields     []Field
	autoSubmit bool
	err        []error
}

// Field represents a typed form field without exposing its concrete type.
type Field interface {
	ID() string
	Label() string
	Focus() tea.Cmd
	Blur()
	IsFocused() bool
	Value() any
	GetValueAsString() string
	Reset()
	Validate() error
	SetError(err error)
	GetError() error
	ResetError()
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
}

func New() *Form {
	return &Form{
		fields:     make([]Field, 0),
		autoSubmit: true,
	}
}

func NewForm(fields ...Field) *Form {
	return New().Add(fields...)
}

// - Methods --//
// Fields
func (m *Form) Add(fields ...Field) *Form {
	m.fields = append(m.fields, fields...)
	return m
}

func (m *Form) Focus() *Form {
	return m.focusOnIndex(0)
}

func (m *Form) Reset() *Form {
	m.resetErrors()

	for _, f := range m.fields {
		f.Reset()
		f.Blur()
	}

	return m.focusOnIndex(0)
}

// Input
func (m *Form) GetValuesAsMappedString() map[string]string {
	formValues := make(map[string]string)

	for _, f := range m.fields {
		formValues[f.ID()] = f.GetValueAsString()
	}

	return formValues
}

// Validation
func (m *Form) Validate() *Form {
	m.resetErrors()

	hasErrors := false
	for _, f := range m.fields {
		if err := f.Validate(); err != nil {
			f.SetError(err)
			hasErrors = true
		}
	}

	if hasErrors {
		m.SetError(logs.Error("error.invalid:form", map[string]any{
			"Error": errors.Join(m.GetFieldsErrors()...),
		}))
	}

	return m
}

func (m *Form) SetError(err error) {
	m.err = append(m.err, err)
}

func (m *Form) GetError() error {
	if len(m.err) > 0 {
		return errors.Join(m.err...)
	}

	return nil
}

func (m *Form) GetFieldsErrors() []error {
	fieldErrors := make([]error, 0, len(m.fields))
	for _, f := range m.fields {
		if err := f.GetError(); err != nil {
			fieldErrors = append(fieldErrors, err)
		}
	}
	return fieldErrors
}

func (m *Form) HasErrors() bool {
	return len(m.err) > 0
}

// Submit
func (m *Form) Submit() *Form {
	m.resetErrors()
	m.Validate()

	return m
}

func (m *Form) EnableAutoSubmit() *Form {
	return m.setAutoSubmit(true)
}

func (m *Form) DisableAutoSubmit() *Form {
	return m.setAutoSubmit(false)
}

func (m *Form) setAutoSubmit(autoSubmit bool) *Form {
	m.autoSubmit = autoSubmit
	return m
}

// - Helper --//
func (m *Form) focusOnIndex(index int) *Form {
	for _, f := range m.fields {
		f.Blur()
	}
	if index >= 0 && index < len(m.fields) {
		m.fields[index].Focus()
	} else if len(m.fields) > 0 {
		m.fields[0].Focus()
	}

	return m
}

func (m *Form) getCurrentField() Field {
	if index := m.getCurrentFieldIndex(); index >= 0 && index < len(m.fields) {
		return m.fields[index]
	}
	return nil
}

func (m *Form) getCurrentFieldIndex() int {
	for i, f := range m.fields {
		if f.IsFocused() {
			return i
		}
	}
	return -1
}

func (m *Form) isLastField() bool {
	currentIndex := m.getCurrentFieldIndex()
	return currentIndex == len(m.fields)-1
}

func (m *Form) resetErrors() {
	m.err = nil
	for _, f := range m.fields {
		f.ResetError()
	}
}

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
			return m, SendMsg[Validate]()
		}
		return m.focusOnIndex(currentIndex + 1), nil
	case PreviousField:
		return m.focusOnIndex(currentIndex - 1), nil
	}

	return m.updateFocusedField(msg)
}
func (m *Form) handleKeyPressMsg(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case sharedtui.KEY_TAB:
		return m, SendMsg[NextField]()
	case sharedtui.KEY_SHIFT_TAB:
		return m, SendMsg[PreviousField]()
	}
	return m.updateFocusedField(msg)
}

func (m *Form) updateFocusedField(msg tea.Msg) (tea.Model, tea.Cmd) {
	index := m.getCurrentFieldIndex()
	if index < 0 {
		return m, nil
	}

	newField, cmd := m.fields[index].Update(msg)
	m.fields[index] = newField.(Field)

	return m, cmd
}
