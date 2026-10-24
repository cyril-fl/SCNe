package footer

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

type plainModel struct{}

func (plainModel) Init() tea.Cmd                       { return nil }
func (plainModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return plainModel{}, nil }
func (plainModel) View() tea.View                      { return tea.NewView("") }

type testModel struct{ plainModel }

func (testModel) GetFooter() []string { return []string{"plugin"} }

func TestFooterView(t *testing.T) {
	f := New("a", "b").Merge(testModel{})

	if got, want := f.View(), "a | b | plugin"; got != want {
		t.Fatalf("View() = %q, want %q", got, want)
	}
	if got := f.GetSize().GetWidth(); got != len("a | b | plugin") {
		t.Fatalf("GetSize().GetWidth() = %d", got)
	}
}

func TestFooterMergeReplacesPreviousPlugins(t *testing.T) {
	f := New("a").Merge(testModel{}).Merge(testModel{})

	if got, want := f.View(), "a | plugin"; got != want {
		t.Fatalf("View() = %q, want %q", got, want)
	}
}

func TestFooterMergeIgnoresModelWithoutFooter(t *testing.T) {
	f := New("a").Merge(plainModel{})

	if got := f.View(); got != "a" {
		t.Fatalf("View() = %q, want a", got)
	}
}
