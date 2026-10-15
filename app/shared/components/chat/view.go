package chat

import (
	"solopg/app/cmdrun/services/game"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/form/field"
	"solopg/app/shared/messages"
	"solopg/app/shared/services/i19n"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/types/primitive"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type Chat struct {
	primitive.Fallible
	primitive.Notifier
	cache

	author   string
	view     viewport.Model
	form     *form.Form
	size     size.Size
	register func(author, message string)
}

type cache struct {
	history []string
	session []string
}

type Template struct {
	Author   string
	History  []string
	Register func(author, message string)
}

func New(params Template) *Chat {
	return &Chat{
		Notifier: primitive.NewNotifier(),

		author:   params.Author,
		form:     makeForm(),
		history:  params.History,
		session:  []string{},
		register: params.Register,
	}
}

// Getters & Setters
func (c *Chat) GetView() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		c.makeViewport(),
		c.spaceInput(),
		c.makeNotification(),
		c.form.View().Content,
		"",
	)
}

func (m Chat) GetFooter() []string {
	return []string{
		i19n.Localize("chat.footer"),
	}
}

func (c *Chat) SetSize(size size.Size) {
	c.size = size
}

// Methods
func (c *Chat) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case form.NextField, form.PreviousField:
		return c.handleInput(msg)
	case form.Validate:
		return c.handleSubmit()

	case messages.Failed:
		return c.handleError(msg)
	case messages.Saved:
		return c.handleSaveNotification()
	case messages.Rolled:
		return c.handleRollNotification()

	case spinner.TickMsg:
		return c.HandleSpinnerTick(msg)
	// case cursor.BlinkMsg:
	// 		return c.handleCursorBlink(msg)

	case tea.KeyMsg:
		return c.handleInput(msg)
	}
	return nil
}

func (c *Chat) Amend(content string) {
	amendChat(c, content)
}

func (c *Chat) makeViewport() string {
	tojoin := []string{}
	row := []string{
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render(makeHistory(c)),

		lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Render(makeWelcomeMessage(c)),

		lipgloss.NewStyle().
			Render(strings.Join(c.session, "\n")),
	}

	for _, r := range row {
		strip := strings.TrimSpace(ansi.Strip(r))

		if strip == "" {
			continue
		}
		tojoin = append(tojoin, r)
	}

	composedView := lipgloss.JoinVertical(
		lipgloss.Left,
		tojoin...,
	)

	c.view.SetContent(composedView)
	c.view.SetWidth(c.size.GetWidth())
	c.view.SetHeight(max(0, c.size.GetHeight()-lipgloss.Height(c.form.View().Content)-2))
	c.view.GotoBottom()

	return c.view.View()
}

func (c *Chat) makeNotification() string {
	if c.IsNotifying() {
		s := lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Render(c.GetSpinnerView())

		m := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render(i19n.Localize(c.GetNotification()))

		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			s,
			m,
		)
	}

	return ""
}

// Helpers
func makeForm() *form.Form {
	return form.
		New(
			field.TextField(field.TextTemplate[string]{
				ID:          "message",
				NoLabel:     true,
				Placeholder: i19n.Localize("chat.placeholder"),
			}),
		).
		SetSendKey(sharedtui.KEY_ENTER).
		Focus()
}

func makeWelcomeMessage(m *Chat) string {
	if len(m.session) == 0 {
		return i19n.Localize("chat.msg:welcome")
	}

	return ""
}

func makeHistory(c *Chat) string {
	return strings.Join(c.history, "\n")
}

func initHistory(engine *game.Engine) []string {
	raw := make([]string, 0, len(engine.State.AdventureLog.Entries))

	for _, entry := range engine.State.AdventureLog.Entries {
		raw = append(raw, entry.String())
	}

	return raw
}

func amendChat(c *Chat, message string) {
	if message == "" {
		return
	}

	c.register(c.author, message)

	authorstyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Render(c.author + ": ")

	cloredmsg := authorstyle + message
	c.session = append(c.session, cloredmsg)
	c.form.Reset()
}

// func (c *Chat) handleCursorBlink(msg cursor.BlinkMsg) tea.Cmd {
// 	var cmd tea.Cmd
// 	_, cmd = c.view.Update(msg)
// 	return cmd
// }

// func getCursor(m model) *tea.Cursor {
// 	// cursor := m.textarea.Cursor()
// 	// cursor.Y += lipgloss.Height(m.viewport.View())

// 	if item := m.getMenuActiveElement(); item != nil && item.IsOpen() {
// 		// cursor = nil
// 	}

// 	return cursor
// }
