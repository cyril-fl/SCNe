package chat

import (
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/cmdrun/tui/shortcut"
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/form"
	"scne/app/shared/components/form/field"
	"scne/app/shared/messages"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/style"
	"scne/app/shared/types/primitive"
	"strings"

	plumprimitive "github.com/cyril-fl/plum/primitive"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	plumtypes "github.com/cyril-fl/plum/types"
)

// Chat shows the session history and a message input with notifications.
type Chat struct {
	plumprimitive.Fallible
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

// Template configures a Chat; Register is called for each sent message.
type Template struct {
	Author   string
	History  []string
	Register func(author, message string)
}

// New returns a chat with a focused message input.
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

// GetView renders the messages, the notification and the input.
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

// GetFooter returns the chat shortcuts.
func (m Chat) GetFooter() []string {
	return []string{
		shortcut.SEND.String(),
	}
}

// SetSize sets the chat size and the input width.
func (c *Chat) SetSize(size size.Size) {
	c.size = size
	c.form.SetWidth(size.GetWidth())
}

// Methods

// Update sends the message on form.Validate, shows save and roll notifications,
// and forwards other messages to the input.
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
	default:
		// paste and cursor blink messages must reach the input
		return c.handleInput(msg)
	}
}

// Amend registers content as a message from the chat author; empty content is ignored.
func (c *Chat) Amend(content string) {
	amendChat(c, content)
}

func (c *Chat) makeViewport() string {
	tojoin := []string{}
	row := []string{
		style.Muted.Render(makeHistory(c)),

		style.Accent.Render(makeWelcomeMessage(c)),

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
		s := style.Accent.Render(c.GetSpinnerView())

		m := style.Muted.Render(i19n.Localize(c.GetNotification()))

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
			field.TextareaField(field.TextareaTemplate{
				ID:          "message",
				NoLabel:     true,
				Placeholder: i19n.Localize("chat.placeholder"),
				Height:      3,
			}),
		).
		SetSendKey(plumtypes.KEY_ENTER).
		SetAutoSubmit(false).
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

func initHistory(gameSession *session.GameSession) []string {
	raw := make([]string, 0, len(gameSession.State.AdventureLog.Entries))

	for _, entry := range gameSession.State.AdventureLog.Entries {
		raw = append(raw, entry.String())
	}

	return raw
}

func amendChat(c *Chat, message string) {
	if message == "" {
		return
	}

	logs.SilentInfo("form.submitted", map[string]any{
		"Form": "chat",
	})

	c.register(c.author, message)

	authorstyle := style.Accent.Render(c.author + ": ")

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
