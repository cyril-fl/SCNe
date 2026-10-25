package statemachine

import (
	"testing"

	"scne/app/cmdrun/types/size"
	"scne/app/shared/services/logs"
	interfass "scne/app/shared/types/interface"

	tea "charm.land/bubbletea/v2"
)

type testView struct {
	content string
	size    size.Size
}

func (v *testView) GetView() string { return v.content }

func (v *testView) Update(tea.Msg) tea.Cmd { return nil }

func (v *testView) SetSize(value size.Size) { v.size = value }

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func TestViewStateTransitions(t *testing.T) {
	chat := &testView{content: "chat"}
	page := &testView{content: "page"}

	state := NewView(TemplateViewState[string]{
		Base:       "chat",
		Transition: map[string][]string{"chat": {"page"}, "page": {"chat"}},
		Views:      map[string]interfass.Viewable{"chat": chat, "page": page},
	})

	if got := state.Current(); got != "chat" {
		t.Fatalf("expected initial state chat, got %q", got)
	}

	if got := state.GetView(); got != "chat" {
		t.Fatalf("expected chat view, got %q", got)
	}

	state.SetSize(size.New(size.Template{Width: 80, Height: 24}))
	state.HandleTransition("page")

	if got := state.Current(); got != "page" {
		t.Fatalf("expected state page, got %q", got)
	}
	if got := state.GetView(); got != "page" {
		t.Fatalf("expected page view, got %q", got)
	}
	if page.size.GetWidth() != 80 || page.size.GetHeight() != 24 {
		t.Fatalf("expected page size 80x24, got %dx%d", page.size.GetWidth(), page.size.GetHeight())
	}
}

func TestViewStateRejectsInvalidTransition(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })

	state := NewView(TemplateViewState[string]{
		Base:       "chat",
		Transition: map[string][]string{"chat": {"page"}},
		Views:      map[string]interfass.Viewable{"chat": &testView{content: "chat"}, "page": &testView{content: "page"}},
	})

	state.HandleTransition("unknown")

	if got := state.Current(); got != "chat" {
		t.Fatalf("expected state to remain chat, got %q", got)
	}
	if !state.HasErr() {
		t.Fatal("expected invalid transition error")
	}
}
