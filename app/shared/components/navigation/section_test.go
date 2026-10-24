package navigation

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

var _ SectionInterface = (*Section)(nil)
var _ FocusableSectionInterface = (*FocusableSection)(nil)

type testNamed string

func (v testNamed) GetName() string { return string(v) }

func TestMenuEnterCallbackReceivesSelectedValue(t *testing.T) {
	var selected testNamed
	menu := NewFocusableSection(TemplateFocusableSection[testNamed]{
		ID:      "test",
		Options: []testNamed{"first", "second"},
		HandleEnterInput: func(value testNamed) tea.Cmd {
			selected = value
			return nil
		},
	})

	if cmd := menu.handlerEnterInput(); cmd != nil {
		t.Fatal("expected callback not to return a command")
	}
	if selected != "first" {
		t.Fatalf("expected selected value %q, got %q", "first", selected)
	}
}
