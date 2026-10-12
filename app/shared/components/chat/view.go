package chat

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/form/field"
	"solopg/app/shared/services/i19n"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/types/primitive"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Chat struct {
	primitive.Fallible
	cache

	author string

	view viewport.Model
	form *form.Form
	size size.Size
}

type cache struct {
	history []string
}

type Template struct {
	Author string
}

func New(params Template) *Chat {
	return &Chat{
		author: params.Author,
		form:   makeForm(),
	}
}

// Getters & Setters
func (c *Chat) GetView() string {
	c.view.SetContent(strings.Join(c.history, "\n"))
	c.view.GotoBottom()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		c.view.View(),
		c.spaceInput(),
		c.form.View().Content,
	)
}

func (c *Chat) GetSize() size.Size {
	return c.size
}

func (c *Chat) SetSize(size size.Size) {
	c.size = size
	c.view.SetWidth(size.GetWidth())
	c.view.SetHeight(max(0, size.GetHeight()-lipgloss.Height(c.form.View().Content)-2))
}

// Methods
func (c *Chat) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case form.NextField, form.PreviousField:
		return c.delegateInputToForm(msg)
	case form.Validate:
		return c.handleFormSubmit()
	case form.Error:
		return c.handleFormError(msg)

	case tea.KeyMsg:
		return c.delegateInputToForm(msg)
	}
	return nil
}

func (c *Chat) spaceInput() string {
	padding := 2
	inputHeight := lipgloss.Height(c.form.View().Content)
	maxHeight := c.size.GetHeight() - inputHeight - padding
	return strings.Repeat("\n", max(0, maxHeight))
}

// Helpers
func makeForm() *form.Form {
	return form.
		New(
			field.TextField(field.TextTemplate[string]{
				ID:          "message",
				Required:    true,
				NoLabel:     true,
				Placeholder: i19n.Localize("gameboard.chat.placeholder"),
			}),
		).
		SetSendKey(sharedtui.KEY_ENTER).
		Focus()
}

func (m Chat) GetFooter() []string {
	return []string{
		"CHAT FOOTER",
	}
}
