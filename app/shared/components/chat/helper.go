package chat

import (
	"scne/app/cmdrun/domain/gameplay/dice"
	"scne/app/shared/messages"
	"scne/app/shared/services/logs"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (c *Chat) handleSubmit() tea.Cmd {
	c.form.Submit()

	if c.form.HasErr() {
		return messages.Fail(c.form.GetErr())
	}

	message, ok := c.form.GetFieldValue("message").(string)
	if !ok {
		return messages.Fail(logs.Error("error.unexpected:value", map[string]any{
			"Subject":  "chat.message",
			"Expected": "string",
			"Value":    c.form.GetFieldValue("message"),
		}))
	}

	amendChat(c, message)
	// Reset only here: Amend also feeds roll results and must keep the text being typed.
	c.form.Reset()
	return nil
}

func (c *Chat) handleError(err messages.Failed) tea.Cmd {
	logs.SilentError("form.error", map[string]any{
		"Form":  "onboardename",
		"Error": err,
	})
	return nil
}

func (c *Chat) handleInput(params tea.Msg) tea.Cmd {
	_, cmd := c.form.Update(params)
	return cmd
}

func (c *Chat) handleSaveNotification() tea.Cmd {
	return c.HandleNotification(1500, "chat.msg:save")
}

func (c *Chat) handleRollNotification() tea.Cmd {
	random := dice.Roll(20)

	return c.HandleNotification(750, "notification:roll"+strconv.Itoa(random))
}

// todo move
func (c *Chat) spaceInput() string {
	padding := 2
	if c.IsNotifying() {
		padding += 1
	}

	inputHeight := lipgloss.Height(c.form.View().Content)
	viewHeight := lipgloss.Height(c.view.View())

	availableHeight := c.size.GetHeight() - viewHeight - inputHeight - padding
	return strings.Repeat("\n", max(0, availableHeight))
}
