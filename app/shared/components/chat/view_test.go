package chat

import (
	"strings"
	"testing"

	"solopg/app/shared/components/form"
	"solopg/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
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
