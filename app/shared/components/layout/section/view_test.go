package section

import (
	"testing"

	"solopg/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type testModel struct {
	content string
	size    *size.Size
}

func (m testModel) Init() tea.Cmd                       { return nil }
func (m testModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m testModel) View() tea.View                      { return tea.NewView(m.content) }
func (m testModel) GetSection() []string                { return []string{m.content} }
func (m testModel) SetSectionSize(s size.Size)          { *m.size = s }

func TestSectionMergeReplacesPreviousPluginView(t *testing.T) {
	s := New("top")
	s.Merge(testModel{content: "first"})
	s.Merge(testModel{content: "second"})

	if got, want := s.View(), lipgloss.JoinVertical(lipgloss.Left, "top", "second"); got != want {
		t.Fatalf("View() = %q, want %q", got, want)
	}
	if got := s.GetSize().GetHeight(); got != 2 {
		t.Fatalf("GetSize().GetHeight() = %d, want 2", got)
	}
}

func TestSectionSetSizeForwardsToModel(t *testing.T) {
	var got size.Size
	want := size.New(size.Template{Width: 10, Height: 5})

	Nil().SetSize(testModel{size: &got}, want)

	if got.GetWidth() != 10 || got.GetHeight() != 5 {
		t.Fatalf("model size = %dx%d, want 10x5", got.GetWidth(), got.GetHeight())
	}
}
