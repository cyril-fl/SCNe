package chat

import (
	"solopg/app/shared/components/form"
	"solopg/app/shared/services/logs"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (c *Chat) handleSubmit() tea.Cmd {
	c.form.Submit()

	if c.form.HasErr() {
		return form.SendErrorMsg(c.form.GetErr())
	}

	logs.SilentInfo("form.submitted", map[string]any{
		"Form": "onboardename",
	})

	message, ok := c.form.GetFieldValue("message").(string)
	if !ok {
		return form.SendErrorMsg(logs.Error("error.unexpected:value", map[string]any{
			"Subject":  "chat.message",
			"Expected": "string",
			"Value":    c.form.GetFieldValue("message"),
		}))
	}

	amendChat(c, message)
	return nil
}

func (c *Chat) handleError(err form.Error) tea.Cmd {
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

func (c *Chat) handleSave() tea.Cmd {
	if c.saving {
		return nil
	}

	c.saving = true
	return tea.Batch(
		c.spinner.Tick,
		tea.Tick(1500*time.Millisecond, func(time.Time) tea.Msg {
			c.saving = false
			return nil
		}),
	)
}

func (c *Chat) handleSpinnerTick(msg tea.Msg) tea.Cmd {
	if !c.saving {
		return nil
	}

	var cmd tea.Cmd
	c.spinner, cmd = c.spinner.Update(msg)
	return cmd
}

func (c *Chat) spaceInput() string {
	padding := 2
	if c.saving {
		padding += 1
	}

	inputHeight := lipgloss.Height(c.form.View().Content)
	viewHeight := lipgloss.Height(c.view.View())

	availableHeight := c.size.GetHeight() - viewHeight - inputHeight - padding
	return strings.Repeat("\n", max(0, availableHeight))
}
