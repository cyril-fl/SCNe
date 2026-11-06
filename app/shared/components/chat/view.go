package chat

import (
	"scne/app/cmdrun/domain/campaign"
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

	author     string
	view       viewport.Model
	form       *form.Form
	size       size.Size
	register   func(author, message string)
	unregister func()
}

type cache struct {
	history []string
	entries primitive.History[entry] // this session
}

type entryKind int

const (
	systemEntry entryKind = iota
	playerEntry
)

// entry is a session message; only player entries can be undone.
type entry struct {
	kind    entryKind
	author  string // registered again on redo
	message string
	line    string
}

// Template configures a Chat; Register is called for each sent message and
// Unregister removes the last registered one when it is undone.
type Template struct {
	Author     string
	History    []string
	Register   func(author, message string)
	Unregister func()
}

// New returns a chat with a focused message input.
func New(params Template) *Chat {
	return &Chat{
		Notifier: primitive.NewNotifier(),

		author:     params.Author,
		form:       makeForm(),
		history:    params.History,
		register:   params.Register,
		unregister: params.Unregister,
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
		shortcut.UNDO.String(),
		shortcut.REDO.String(),
	}
}

// SetSize sets the chat size and the input width.
func (c *Chat) SetSize(size size.Size) {
	c.size = size
	c.form.SetWidth(size.GetWidth())
}

// Methods

// Update sends the message on form.Validate, undoes or redoes player messages
// on ctrl+z / ctrl+y, shows save and roll notifications, and forwards other
// messages to the input.
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
		switch plumtypes.Key(msg.String()) {
		case plumtypes.SHORT_CTRL_Z:
			c.undo()
			return nil
		case plumtypes.SHORT_CTRL_Y:
			c.redo()
			return nil
		}
		return c.handleInput(msg)
	default:
		// paste and cursor blink messages must reach the input
		return c.handleInput(msg)
	}
}

// Amend registers content as a system entry under the chat author, so undo
// stops on it; empty content is ignored.
func (c *Chat) Amend(content string) {
	amendChat(c, content, systemEntry)
}

func (c *Chat) makeViewport() string {
	tojoin := []string{}
	row := []string{
		makeHistory(c),

		style.Accent().Render(makeWelcomeMessage(c)),

		lipgloss.NewStyle().
			Render(strings.Join(sessionLines(c), "\n")),
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
		s := style.Accent().Render(c.GetSpinnerView())

		m := style.Muted().Render(i19n.Localize(c.GetNotification()))

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
	if len(m.entries.Passed()) == 0 {
		return i19n.Localize("chat.msg:welcome")
	}

	return ""
}

func makeHistory(c *Chat) string {
	lines := make([]string, 0, len(c.history))
	for _, line := range c.history {
		lines = append(lines, styleRoll(line, style.Muted()))
	}
	return strings.Join(lines, "\n")
}

const rollArrow = " → "

// isRoll reports whether message is a roll result; rolls are shown without author.
func isRoll(message string) bool {
	return strings.Contains(message, rollArrow)
}

// styleRoll mutes a roll's context and highlights its result after the last
// " → ", with a trailing "✦" critical mark in the error color. Other lines get
// plain.
// shortcut: any line containing " → " is styled as a roll, store an entry kind
// in the journal if players start typing arrows.
func styleRoll(line string, plain lipgloss.Style) string {
	i := strings.LastIndex(line, rollArrow)
	if i < 0 {
		return plain.Render(line)
	}

	result, critical, isCritical := strings.Cut(line[i+len(rollArrow):], " ✦ ")
	styled := style.Muted().Render(line[:i]+" →") + " " + style.Accent().Bold(true).Render(result)
	if isCritical {
		styled += " " + style.Error().Render("✦ "+critical)
	}
	return styled
}

func initHistory(gameSession *session.GameSession) []string {
	raw := make([]string, 0, len(gameSession.State.AdventureLog.Entries))

	for _, entry := range gameSession.State.AdventureLog.Entries {
		raw = append(raw, entry.String())
	}

	return raw
}

func amendChat(c *Chat, message string, kind entryKind) {
	if message == "" {
		return
	}

	logs.SilentInfo("form.submitted", map[string]any{
		"Form": "chat",
	})

	if isRoll(message) {
		c.register("", message)
		c.entries.Push(entry{kind: systemEntry, message: message, line: styleRoll(message, lipgloss.NewStyle())})
		return
	}

	c.register(c.author, message)

	authorstyle := style.Accent().Render(campaign.DisplayAuthor(c.author) + ": ")
	c.entries.Push(entry{kind: kind, author: c.author, message: message, line: authorstyle + message})
}

func sessionLines(c *Chat) []string {
	passed := c.entries.Passed()
	lines := make([]string, 0, len(passed))
	for _, e := range passed {
		lines = append(lines, e.line)
	}
	return lines
}

// undo removes the last player entry; a system entry stops it so rolls cannot
// be undone and rerolled.
func (c *Chat) undo() {
	if _, ok := c.entries.Undo(isPlayer); ok && c.unregister != nil {
		c.unregister()
	}
}

// redo registers the last undone entry again.
func (c *Chat) redo() {
	if e, ok := c.entries.Redo(); ok {
		c.register(e.author, e.message)
	}
}

func isPlayer(e entry) bool {
	return e.kind == playerEntry
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
