package navigation

import (
	"strings"
	"testing"

	"solopg/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
)

var _ SectionInterface = (*Section)(nil)
var _ FocusableSectionInterface = (*FocusableSection)(nil)

type testNamed string

func (v testNamed) GetName() string { return string(v) }

func TestFocusableSectionShowsOptionsOnlyWhenFocused(t *testing.T) {
	menu := NewFocusableSection(TemplateFocusableSection[testNamed]{
		ID:      "test",
		Options: []testNamed{"first"},
	})
	menu.SetSize(size.New(size.Template{Width: 40, Height: 20}))

	if view := menu.GetView(); strings.Contains(view, "first") || !strings.Contains(view, "TEST") {
		t.Fatalf("expected unfocused view to show the title only, got %q", view)
	}

	menu.SetFocus(true)
	if view := menu.GetView(); !strings.Contains(view, "first") {
		t.Fatalf("expected focused view to list the options, got %q", view)
	}
}

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
