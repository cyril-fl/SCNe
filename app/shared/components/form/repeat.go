package form

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

// - Repeat - //

// repeatField is a Field duplicating a group of fields: it manages its copies
// in an inner Form, the way Form manages its fields.
type repeatField struct {
	id    string
	label string
	make  func(index int) []Field
	count int
	group *Form
	err   error
}

type RepeatTemplate struct {
	ID    string
	Label string
	// Count is the number of groups to start with, at least one.
	Count int
	// Make builds the fields of the group at index.
	Make func(index int) []Field
}

func Repeat(template RepeatTemplate) *repeatField {
	r := &repeatField{
		id:    template.ID,
		label: template.Label,
		make:  template.Make,
		group: new(),
	}

	for range max(1, template.Count) {
		r.addGroup()
	}

	return r
}

// - Field Implementation - //
func (r *repeatField) ID() string {
	return r.id
}

func (r *repeatField) Label() string {
	return r.label
}

func (r *repeatField) Focus() tea.Cmd {
	r.group.Focus()
	return nil
}

func (r *repeatField) Blur() {
	for _, f := range r.group.fields {
		f.Blur()
	}
}

func (r *repeatField) IsFocused() bool {
	return r.group.getCurrentFieldIndex() >= 0
}

func (r *repeatField) Value() any {
	return r.group.GetValuesAsMappedString()
}

// GetValueAsString is empty: the form reads each repeated field through Value.
func (r *repeatField) GetValueAsString() string {
	return ""
}

func (r *repeatField) Reset() {
	r.group.Reset()
	r.Blur()
}

func (r *repeatField) Validate() error {
	r.group.Validate()
	// The inner form error would repeat the field errors in the view.
	r.group.ClearErr()

	return errors.Join(r.group.getFieldsErrors()...)
}

func (r *repeatField) SetError(err error) {
	r.err = err
}

func (r *repeatField) GetError() error {
	return r.err
}

func (r *repeatField) ResetError() {
	r.err = nil
	r.group.handleErrorsReset()
}

func (r *repeatField) Init() tea.Cmd {
	return r.group.Init()
}

func (r *repeatField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch plumtypes.Key(key.String()) {
		// Terminals report a shifted "+" either way.
		case "ctrl++", "ctrl+=":
			first := len(r.group.fields)
			r.addGroup()
			r.group.focusOnIndex(first)
			return r, nil
		}
	}

	_, cmd := r.group.updateFocusedField(msg)
	return r, cmd
}

func (r *repeatField) View() tea.View {
	return r.group.View()
}

func (r *repeatField) SetWidth(width int) {
	r.group.SetWidth(width)
}

// - Navigation - //

// next moves the focus to the next repeated field, false when already on the last one.
func (r *repeatField) next() bool {
	index := r.group.getCurrentFieldIndex()
	if index < 0 || index >= len(r.group.fields)-1 {
		return false
	}
	r.group.focusOnIndex(index + 1)
	return true
}

// previous moves the focus to the previous repeated field, false when already on the first one.
func (r *repeatField) previous() bool {
	index := r.group.getCurrentFieldIndex()
	if index <= 0 {
		return false
	}
	r.group.focusOnIndex(index - 1)
	return true
}

func (r *repeatField) focusLast() {
	r.group.focusOnIndex(len(r.group.fields) - 1)
}

// Helpers
func (r *repeatField) addGroup() {
	r.group.Add(r.make(r.count)...)
	r.count++
}
