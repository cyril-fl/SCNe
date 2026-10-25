package menuaction

import (
	"scne/app/cmdrun/tui/shortcut"
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/form"
	"scne/app/shared/messages"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

// PageTemplate configures an action page.
type PageTemplate struct {
	// Label is the footer label of the submit shortcut.
	Label   string
	NewForm func() *form.Form
	// Submit runs the action from the validated form.
	Submit func(values *form.Form) tea.Cmd
}

// actionPage is a formpage.Page without the list: one form that runs an action on submit.
// The form is rebuilt on each opening so its options follow the session.
type actionPage struct {
	label   string
	newForm func() *form.Form
	submit  func(values *form.Form) tea.Cmd

	form *form.Form
}

func newPage(template PageTemplate) *actionPage {
	return &actionPage{
		label:   template.Label,
		newForm: template.NewForm,
		submit:  template.Submit,
	}
}

// Getters & Setters

// GetView renders the form.
func (p *actionPage) GetView() string {
	if p.form == nil {
		return ""
	}
	return p.form.View().Content
}

// GetFooter returns the submit shortcut.
func (p *actionPage) GetFooter() []string {
	return []string{shortcut.NewShortcut(shortcut.Template{
		Command: shortcut.CtrlEnter,
		Label:   p.label,
	}).String()}
}

// SetSize builds the form when the page is entered, then sets its width.
func (p *actionPage) SetSize(s size.Size) {
	if p.form == nil {
		p.form = p.newForm()
	}
	p.form.SetWidth(s.GetWidth())
}

// Methods

// Update submits on form.Validate, closes on esc, and routes other messages to the form.
func (p *actionPage) Update(msg tea.Msg) tea.Cmd {
	if p.form == nil {
		return nil
	}

	switch msg := msg.(type) {
	case form.Validate:
		return p.handleValidate()
	case tea.KeyPressMsg:
		return p.handleKeyPress(msg)
	default:
		_, cmd := p.form.Update(msg)
		return cmd
	}
}

// Helpers

func (p *actionPage) handleKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	switch plumtypes.Key(msg.String()) {
	case plumtypes.KEY_ESC:
		return p.close()
	case plumtypes.CMD_CTRL_ENTER:
		return form.Send[form.Validate]()
	default:
		_, cmd := p.form.Update(msg)
		return cmd
	}
}

func (p *actionPage) handleValidate() tea.Cmd {
	p.form.Submit()
	if p.form.HasErr() {
		return messages.Fail(p.form.GetErr())
	}

	action := p.submit(p.form)
	closing := p.close()

	return tea.Sequence(closing, action)
}

func (p *actionPage) close() tea.Cmd {
	p.form = nil
	return form.Send[form.Close]()
}
