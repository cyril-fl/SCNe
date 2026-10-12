package sidepanel

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

type testModel struct {
	content string
}

func (m testModel) Init() tea.Cmd { return nil }

func (m testModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m testModel) View() tea.View { return tea.NewView(m.content) }

func (m testModel) GetSidepanel() []string { return []string{m.content} }

func TestMergeReplacesPreviousPluginView(t *testing.T) {
	panel := Nil()
	panel.Merge(testModel{content: "first"})
	panel.Merge(testModel{content: "second"})

	if got, want := panel.View(), "second"; got != want {
		t.Fatalf("expected sidepanel %q, got %q", want, got)
	}
}
