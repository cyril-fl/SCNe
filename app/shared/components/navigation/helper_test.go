package navigation

import (
	"testing"

	"solopg/app/cmdrun/types/size"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type testSection struct{ view string }

func (s testSection) ID() string                   { return "" }
func (s testSection) SetFocus(bool)                {}
func (s testSection) GetList() *list.Model         { return nil }
func (s testSection) SetList(*list.Model)          {}
func (s testSection) SetSize(size.Size)            {}
func (s testSection) GetView() string              { return s.view }
func (s testSection) GetFooter() []string          { return nil }
func (s testSection) HandleUpdate(tea.Msg) tea.Cmd { return nil }

func TestMakeNavigationContentTracksBlockOffsets(t *testing.T) {
	content, offsets := makeNavigationContent([]FocusableSectionInterface{testSection{view: "a"}, testSection{view: "bb"}})
	if content != "a\nbb" || len(offsets) != 2 || offsets[0] != 0 || offsets[1] != 1 {
		t.Fatalf("got content %q and offsets %v", content, offsets)
	}
}
