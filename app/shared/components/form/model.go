package form

import (
	"errors"
	plumprimitive "github.com/cyril-fl/plum/primitive"
	"maps"
	"solopg/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

// - Form - //

// Form is a tea model holding an ordered list of fields with a single focused one.
type Form struct {
	plumprimitive.Fallible
	sendKey    plumtypes.Key
	fields     []Field
	autoSubmit bool
}

// Makefunc builds a form prefilled from a value map keyed by field ID.
type Makefunc = func(map[string]string) *Form

// New returns a form with the given fields and auto-submit enabled.
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

// GetFieldValue returns the value of the field with fieldID, or nil when absent.
func (m *Form) GetFieldValue(fieldID string) any {
	for _, f := range m.fields {
		if f.ID() == fieldID {
			return f.Value()
		}
	}

	return nil
}

// GetValuesAsMappedString returns every field value as a string keyed by field ID.
// Repeat fields are flattened into their inner field IDs.
func (m *Form) GetValuesAsMappedString() map[string]string {
	formValues := make(map[string]string)

	for _, f := range m.fields {
		if r, ok := f.(*repeatField); ok {
			maps.Copy(formValues, r.group.GetValuesAsMappedString())
			continue
		}
		formValues[f.ID()] = f.GetValueAsString()
	}

	return formValues
}

// SetSendKey sets the key that validates the form.
func (m *Form) SetSendKey(key plumtypes.Key) *Form {
	m.sendKey = key
	return m
}

// SetWidth resizes every field that supports a width.
func (m *Form) SetWidth(width int) *Form {
	for _, f := range m.fields {
		if resizable, ok := f.(interface{ SetWidth(int) }); ok {
			resizable.SetWidth(width)
		}
	}
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

// CanRepeat reports whether the focused field is a Repeat, which then accepts alt+"+".
func (m *Form) CanRepeat() bool {
	index := m.getCurrentFieldIndex()
	if index < 0 {
		return false
	}
	_, ok := m.fields[index].(*repeatField)
	return ok
}

// SetAutoSubmit sets whether leaving the last field validates the form.
func (m *Form) SetAutoSubmit(autoSubmit bool) *Form {
	m.autoSubmit = autoSubmit
	return m
}

// Methods
// - Tea Model Implementation --//

// Init returns the batched Init commands of all fields.
func (m *Form) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.fields))
	for _, f := range m.fields {
		cmds = append(cmds, f.Init())
	}
	return tea.Batch(cmds...)
}

// Update handles field navigation and forwards other messages to the focused field.
// It does nothing when no field is focused.
func (m *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentIndex := m.getCurrentFieldIndex()
	if currentIndex == -1 {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPressMsg(msg)
	case NextField:
		if r, ok := m.fields[currentIndex].(*repeatField); ok && r.next() {
			return m, nil
		}
		if m.isLastField() && m.autoSubmit {
			return m, Send[Validate]()
		}
		return m.focusOnIndex(currentIndex + 1), nil
	case PreviousField:
		if r, ok := m.fields[currentIndex].(*repeatField); ok && r.previous() {
			return m, nil
		}
		m.focusOnIndex(currentIndex - 1)
		if r, ok := m.fields[m.getCurrentFieldIndex()].(*repeatField); ok {
			r.focusLast()
		}
		return m, nil

	}

	return m.updateFocusedField(msg)
}

// - Fields - //

// Add appends fields to the form.
func (m *Form) Add(fields ...Field) *Form {
	m.fields = append(m.fields, fields...)
	return m
}

// Focus focuses the first field.
func (m *Form) Focus() *Form {
	return m.focusOnIndex(0)
}

// Submit clears previous errors and validates the form.
func (m *Form) Submit() *Form {
	m.handleErrorsReset()
	m.Validate()

	return m
}

// Validate validates every field, storing field errors and a joined form error.
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

// Reset clears errors, resets every field and focuses the first one.
func (m *Form) Reset() *Form {
	m.handleErrorsReset()

	for _, f := range m.fields {
		f.Reset()
		f.Blur()
	}

	return m.focusOnIndex(0)
}

// - Submit - //

// EnableAutoSubmit makes leaving the last field validate the form.
func (m *Form) EnableAutoSubmit() *Form {
	return m.SetAutoSubmit(true)
}

// DisableAutoSubmit stops leaving the last field from validating the form.
func (m *Form) DisableAutoSubmit() *Form {
	return m.SetAutoSubmit(false)
}
