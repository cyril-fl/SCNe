package header

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

type testModel struct{}

func (testModel) Init() tea.Cmd                       { return nil }
func (testModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return testModel{}, nil }
func (testModel) View() tea.View                      { return tea.NewView("") }
func (testModel) GetHeader() Template {
	return Template{Title: "plugin.title", Subtitle: "plugin.subtitle"}
}

func TestHeaderViewWithoutTitleIsEmpty(t *testing.T) {
	if got := Nil().View(); got != "" {
		t.Fatalf("View() = %q, want empty", got)
	}
}

func TestHeaderSetHeaderRendersTitleAndSubtitle(t *testing.T) {
	h := New(Template{Title: "initial.title"})
	h.SetHeader(testModel{})

	view := h.View()
	if !strings.Contains(view, "plugin.title") || !strings.Contains(view, "plugin.subtitle") {
		t.Fatalf("View() = %q, want title and subtitle", view)
	}
	if strings.Contains(view, "initial.title") {
		t.Fatal("SetHeader() did not replace the title")
	}
}
