package field

import (
	"fmt"
	"solopg/app/cmdrun/tui/models"
	"solopg/app/cmdrun/types/direction"
	"solopg/app/shared/components/form"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/utils/transform"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// - Field - //
type selectField[T any] struct {
	field[T]
	options list.Model
}

type SelectTemplate[T any] struct {
	ID           string
	Label        string
	Options      []models.Item[T]
	Defaultvalue T
	Validator    func(T) error
	Required     bool
}

func SelectField[T any](template SelectTemplate[T]) *selectField[T] {
	items := make([]list.Item, 0, len(template.Options))
	for _, option := range template.Options {
		items = append(items, option)
	}
	// FIXME LOW Mettre des taille c'est juste un fix tempraire, ca le le rend pas "flex"
	options := list.New(items, list.NewDefaultDelegate(), 15, 10)
	models.ConfigureList(&options)
	models.SetListFocus(&options, false)

	return &selectField[T]{
		field:   newField(template.ID, template.Label, SELECT, template.Defaultvalue, template.Validator, template.Required),
		options: options,
	}
}

// - Field Implementation - //
func (f *selectField[T]) Focus() tea.Cmd {
	f.focus = true
	models.SetListFocus(&f.options, true)
	return nil
}

func (f *selectField[T]) Blur() {
	f.focus = false
	models.SetListFocus(&f.options, false)
}

func (f *selectField[T]) Value() any {
	return f.options.SelectedItem().(models.Item[T]).Value()
}

func (f *selectField[T]) GetValueAsString() string {
	return transform.ParseAsString(f.Value())
}

func (f *selectField[T]) Reset() {
	f.ResetError()
	f.Blur()
	f.options.Select(0)
}

func (f *selectField[T]) Validate() error {
	if err := f.testRequireness(); err != nil {
		return err
	}

	if v, ok := f.Value().(T); ok {
		return f.validate(v)
	}

	return logs.Error("error.unexpected:value", map[string]any{
		"Subject":  "field",
		"Expected": fmt.Sprintf("%T", f.defaultvalue),
		"Value":    f.Value(),
	})
}

func (f *selectField[T]) testRequireness() error {
	if f.required && f.Value() == nil {
		return logs.Error("error.required", map[string]any{
			"Subject": f.Label(),
			"Value":   f.ID(),
		})
	}
	return nil
}

// - Tea Model Implementation - //
func (m *selectField[T]) Init() tea.Cmd {
	return nil
}

func (m *selectField[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.options.SetSize(msg.Width, msg.Height)
	case tea.KeyPressMsg:
		switch msg.String() {
		case sharedtui.KEY_UP, sharedtui.KEY_DOWN:
			return m.handleDirection(msg)
		case sharedtui.KEY_ENTER:
			return m, tea.Sequence(
				form.SendMsg[form.NextField](),
				form.SendMsg[form.Scroll](),
			)
		}
	default:
		newOptions, cmd := m.options.Update(msg)
		m.options = newOptions
		return m, cmd
	}
	return m, nil
}

func (m *selectField[T]) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Render(m.Label())

	options := m.options.View()

	err := ""
	if m.GetError() != nil {
		err = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Render(m.GetError().Error())
	}

	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			err,
			options,
		),
	)
}

// - Helper - //
func (m *selectField[T]) handleDirection(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	updatedList, newdirection := direction.GetListDirection(&m.options, msg)
	m.options = updatedList
	switch newdirection {
	case direction.NEXT:
		return m, tea.Sequence(
			form.SendMsg[form.NextField](),
			form.SendMsg[form.Scroll](),
		)
	case direction.PREVIOUS:
		return m, tea.Sequence(
			form.SendMsg[form.PreviousField](),
			form.SendMsg[form.Scroll](),
		)
	}
	return m, nil
}
