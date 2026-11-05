package chat

import (
	"strings"
	"testing"

	"scne/app/shared/components/form"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/style"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func TestChatSubmitRegistersMessage(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })

	var author, message string
	c := New(Template{
		Author:   "Ada",
		Register: func(a, m string) { author, message = a, m },
	})

	c.Update(tea.KeyPressMsg{Code: 'h', Text: "h"})
	c.Update(tea.KeyPressMsg{Code: 'i', Text: "i"})
	c.Update(form.Validate{})

	if author != "Ada" || message != "hi" {
		t.Fatalf("Register received %q, %q; want Ada, hi", author, message)
	}
	if len(c.session) != 1 || !strings.Contains(c.session[0], "hi") {
		t.Fatalf("session = %v, want the submitted message", c.session)
	}
	if got := c.form.GetFieldValue("message"); got != "" {
		t.Fatalf("message after submit = %q, want empty", got)
	}
}

func TestChatSubmitIgnoresEmptyMessage(t *testing.T) {
	called := false
	c := New(Template{Register: func(string, string) { called = true }})

	c.Update(form.Validate{})

	if called || len(c.session) != 0 {
		t.Fatal("empty message was registered")
	}
}

func TestChatAmendKeepsTypedMessage(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })

	c := New(Template{Register: func(string, string) {}})

	c.Update(tea.KeyPressMsg{Code: 'h', Text: "h"})
	c.Amend("rolled 4")

	if got := c.form.GetFieldValue("message"); got != "h" {
		t.Fatalf("message after Amend = %q, want h", got)
	}
}

func TestStyleRollHighlightsResultAndCritical(t *testing.T) {
	plain := lipgloss.NewStyle()
	tests := []struct {
		line   string
		result string
		mark   string
	}{
		{line: "hello", result: ""},
		{line: "Probable · 19 → Oui", result: "Oui"},
		{line: "Probable · 20 → Oui ✦ critique", result: "Oui", mark: "✦ critique"},
	}

	for _, tt := range tests {
		got := styleRoll(tt.line, plain)
		if ansi.Strip(got) != tt.line {
			t.Fatalf("styleRoll(%q) text = %q, want unchanged", tt.line, ansi.Strip(got))
		}
		if tt.result != "" && !strings.Contains(got, style.Accent().Bold(true).Render(tt.result)) {
			t.Fatalf("styleRoll(%q) = %q, want highlighted %q", tt.line, got, tt.result)
		}
		if tt.mark != "" && !strings.Contains(got, style.Error().Render(tt.mark)) {
			t.Fatalf("styleRoll(%q) = %q, want critical mark %q", tt.line, got, tt.mark)
		}
	}
}
