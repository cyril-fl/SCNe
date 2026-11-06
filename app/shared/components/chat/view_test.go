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
	if len(c.entries.Passed()) != 1 || !strings.Contains(c.entries.Passed()[0].line, "hi") {
		t.Fatalf("passed = %v, want the submitted message", c.entries.Passed())
	}
	if got := c.form.GetFieldValue("message"); got != "" {
		t.Fatalf("message after submit = %q, want empty", got)
	}
}

func TestChatSubmitIgnoresEmptyMessage(t *testing.T) {
	called := false
	c := New(Template{Register: func(string, string) { called = true }})

	c.Update(form.Validate{})

	if called || len(c.entries.Passed()) != 0 {
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

func TestChatAmendRegistersRollWithoutAuthor(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })

	author := "unset"
	c := New(Template{
		Author:   "Ada",
		Register: func(a, _ string) { author = a },
	})

	c.Amend("d20 → 14")
	if author != "" {
		t.Fatalf("Register author = %q, want empty for a roll", author)
	}
	if got := ansi.Strip(c.entries.Passed()[0].line); got != "d20 → 14" {
		t.Fatalf("session line = %q, want the roll without author", got)
	}
}

func send(c *Chat, text string) {
	for _, r := range text {
		c.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	c.Update(form.Validate{})
}

var (
	ctrlZ = tea.KeyPressMsg{Code: 'z', Mod: tea.ModCtrl}
	ctrlY = tea.KeyPressMsg{Code: 'y', Mod: tea.ModCtrl}
)

func TestChatUndoRedoPlayerMessage(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })

	var journal []string
	c := New(Template{
		Author:     "Ada",
		Register:   func(_, m string) { journal = append(journal, m) },
		Unregister: func() { journal = journal[:len(journal)-1] },
	})

	send(c, "a")
	send(c, "b")
	c.Update(ctrlZ)
	c.Update(ctrlZ)
	if len(c.entries.Passed()) != 0 || len(journal) != 0 {
		t.Fatalf("after undo passed = %v, journal = %v, want both empty", c.entries.Passed(), journal)
	}

	c.Update(ctrlY)
	if len(c.entries.Passed()) != 1 || c.entries.Passed()[0].message != "a" || strings.Join(journal, ",") != "a" {
		t.Fatalf("after redo passed = %v, journal = %v, want a", c.entries.Passed(), journal)
	}
}

func TestChatUndoStopsOnSystemEntry(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })

	unregistered := 0
	c := New(Template{
		Author:     "Ada",
		Register:   func(string, string) {},
		Unregister: func() { unregistered++ },
	})

	send(c, "a")
	c.Amend("d20 → 14")
	c.Update(ctrlZ)

	if len(c.entries.Passed()) != 2 || unregistered != 0 {
		t.Fatalf("undo went past the roll: passed = %v, unregistered = %d", c.entries.Passed(), unregistered)
	}
}

func TestChatNewEntryDumpsFuture(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })

	c := New(Template{Author: "Ada", Register: func(string, string) {}, Unregister: func() {}})

	send(c, "a")
	c.Update(ctrlZ)
	send(c, "b")
	c.Update(ctrlY)

	if len(c.entries.Passed()) != 1 || c.entries.Passed()[0].message != "b" {
		t.Fatalf("passed = %v, want only b after redo with an empty future", c.entries.Passed())
	}
}
