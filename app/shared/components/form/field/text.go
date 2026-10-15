package field

import (
	"fmt"
	"solopg/app/shared/components/form"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/utils/transform"
	"strings"

	"strconv"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// - Field - //
type textField[T stringOrInt] struct {
	field[T]
	input textinput.Model
}

type TextTemplate[T stringOrInt] struct {
	Template[T]
	Input       textinput.Model
	Placeholder string
}

func TextField[T stringOrInt](template TextTemplate[T]) *textField[T] {
	input := makeInput(template.Placeholder)
	setValue(&input, template.Defaultvalue)

	return &textField[T]{
		field: newField(
			template.ID,
			template.Label,
			INPUT,
			template.Defaultvalue,
			template.Validator,
			template.Required,
			template.NoLabel,
		),
		input: input,
	}
}

func (m *textField[T]) View() tea.View {
	err := ""
	if m.GetError() != nil {
		err = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Render(m.GetError().Error())
	}

	if m.nolabel {
		return m.getNoLabelView(m.input.View(), err)
	} else {
		return m.getLabelView(m.input.View(), err)

	}
}

// Getters & Setters
func (f *textField[T]) GetValueAsString() string {
	return transform.ParseAsString(f.Value())
}

func (m *textField[T]) getNoLabelView(input, err string) tea.View {
	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			err,
			input,
		),
	)
}

func (m *textField[T]) getLabelView(input, err string) tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Render(m.Label())

	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			err,
			input,
		),
	)
}

func (m *textField[T]) SetWidth(size int) {
	m.input.SetWidth(size)
}

// Methods
// - Tea Model Implementation - //
func (m *textField[T]) Init() tea.Cmd {
	return textinput.Blink
}

func (m *textField[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch sharedtui.EventKey(msg.String()) {
		case sharedtui.KEY_ENTER:
			return m, tea.Sequence(
				form.Send[form.NextField](),
				form.Send[form.Scroll](),
			)
		case sharedtui.KEY_UP:
			return m, tea.Sequence(
				form.Send[form.PreviousField](),
				form.Send[form.Scroll](),
			)
		case sharedtui.KEY_DOWN:
			return m, tea.Sequence(
				form.Send[form.NextField](),
				form.Send[form.Scroll](),
			)
		default:
			newInput, cmd := m.input.Update(msg)
			m.input = newInput
			return m, cmd
		}
	}

	return m, nil
}

// - Field Implementation - //
func (f *textField[T]) Focus() tea.Cmd {
	f.focus = true
	return f.input.Focus()
}

func (f *textField[T]) Blur() {
	f.focus = false
	f.input.Blur()
}

func (f *textField[T]) Value() any {
	return f.input.Value()
}

func (f *textField[T]) Reset() {
	f.ResetError()
	f.Blur()
	setValue(&f.input, f.defaultvalue)
}

func (f *textField[T]) Validate() error {
	if err := f.testRequireness(); err != nil {
		return err
	}

	if v, ok := f.Value().(T); ok {
		return f.validate(v)
	}

	return logs.Error("error.unexpected:value", map[string]any{
		"Subject":  f.Label(),
		"Expected": fmt.Sprintf("%T", f.defaultvalue),
		"Value":    f.Value(),
	})
}

func (f *textField[T]) testRequireness() error {
	if f.required && strings.TrimSpace(f.Value().(string)) == "" {
		return logs.Error("error.required", map[string]any{
			"Subject":  f.ID(),
			"Property": f.Label(),
		})
	}
	return nil
}

// - Helper - //
type stringOrInt interface {
	~string | ~int
}

func makeInput(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	w := lipgloss.Width(placeholder)
	input.SetWidth(w)
	return input
}

func setValue[T stringOrInt](input *textinput.Model, value T) {
	switch v := any(value).(type) {
	case string:
		input.SetValue(v)
	case int:
		input.SetValue(strconv.Itoa(v))
	}
}
