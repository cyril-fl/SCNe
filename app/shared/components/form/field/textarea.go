package field

import (
	"scne/app/shared/components/form"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/style"
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

// - Field - //
type textareaField struct {
	field[string]
	input textarea.Model
}

// TextareaTemplate configures a TextareaField; Height is at least one line.
type TextareaTemplate struct {
	Template[string]
	Placeholder string
	Height      int
}

// TextareaField returns a multi-line input field; alt+enter inserts a new line.
func TextareaField(template TextareaTemplate) *textareaField {
	input := makeTextarea(template.Placeholder, template.Height)
	input.SetValue(template.Defaultvalue)

	return &textareaField{
		field: newField(
			template.ID,
			template.Label,
			TEXTAREA,
			template.Defaultvalue,
			template.Validator,
			template.Required,
			template.NoLabel,
		),
		input: input,
	}
}

func (m *textareaField) View() tea.View {
	err := ""
	if m.GetError() != nil {
		err = style.Error.Render(m.GetError().Error())
	}

	if m.nolabel {
		return m.getNoLabelView(m.input.View(), err)
	} else {
		return m.getLabelView(m.input.View(), err)
	}
}

// Getters & Setters
func (f *textareaField) GetValueAsString() string {
	return f.input.Value()
}

func (m *textareaField) getNoLabelView(input, err string) tea.View {
	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			err,
			input,
		),
	)
}

func (m *textareaField) getLabelView(input, err string) tea.View {
	title := style.Title.Render(m.Label())

	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			err,
			input,
		),
	)
}

func (m *textareaField) SetWidth(size int) {
	m.input.SetWidth(size)
}

// Methods
// - Tea Model Implementation - //
func (m *textareaField) Init() tea.Cmd {
	return textarea.Blink
}

func (m *textareaField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch plumtypes.Key(msg.String()) {
		case plumtypes.KEY_ENTER:
			return m, tea.Sequence(
				form.Send[form.NextField](),
				form.Send[form.Scroll](),
			)
		case plumtypes.CMD_ALT_ENTER:
			if strings.TrimSpace(m.input.Value()) == "" {
				return m, nil
			}
			newInput, cmd := m.input.Update(msg)
			m.input = newInput
			return m, cmd
		default:
			newInput, cmd := m.input.Update(msg)
			m.input = newInput
			return m, cmd
		}
	}

	// paste and cursor blink messages are not key presses
	newInput, cmd := m.input.Update(msg)
	m.input = newInput
	return m, cmd
}

// - Field Implementation - //
func (f *textareaField) Focus() tea.Cmd {
	f.focus = true
	return f.input.Focus()
}

func (f *textareaField) Blur() {
	f.focus = false
	f.input.Blur()
}

func (f *textareaField) Value() any {
	return f.input.Value()
}

func (f *textareaField) Reset() {
	f.ResetError()
	f.Blur()
	f.input.SetValue(f.defaultvalue)
}

func (f *textareaField) Validate() error {
	if err := f.testRequireness(); err != nil {
		return err
	}

	return f.validate(f.input.Value())
}

func (f *textareaField) testRequireness() error {
	if f.required && strings.TrimSpace(f.input.Value()) == "" {
		return logs.Error("error.required", map[string]any{
			"Subject":  f.ID(),
			"Property": f.Label(),
		})
	}
	return nil
}

// - Helper - //
// makeTextarea mirrors the bubbletea chat example: no line numbers, no cursor line highlight, enter is left to the form.
func makeTextarea(placeholder string, height int) textarea.Model {
	input := textarea.New()
	input.Placeholder = placeholder
	input.ShowLineNumbers = false
	input.KeyMap.InsertNewline.SetKeys(string(plumtypes.CMD_ALT_ENTER))
	input.KeyMap.WordForward.SetKeys("alt+right", "alt+f", "ctrl+right")
	input.KeyMap.WordBackward.SetKeys("alt+left", "alt+b", "ctrl+left")

	styles := input.Styles()
	styles.Focused.CursorLine = lipgloss.NewStyle()
	input.SetStyles(styles)

	input.SetHeight(max(1, height))
	return input
}
