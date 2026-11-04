package navigation

import (
	"testing"

	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/models"
	"scne/app/shared/services/i19n"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
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
	if content != "a\n\nbb" || len(offsets) != 2 || offsets[0] != 0 || offsets[1] != 2 {
		t.Fatalf("got content %q and offsets %v", content, offsets)
	}
}

func TestMakeSectionModelCapitalizesOptionNames(t *testing.T) {
	translator := i19n.New()
	translator.SetLocale(i18n.NewLocalizer(i18n.NewBundle(language.English), "en"))
	i19n.SetCache(translator)
	t.Cleanup(func() { i19n.SetCache(nil) })

	model := makeSectionModel(TemplateFocusableSection[testNamed]{Options: []testNamed{"fight"}})

	if got := model.Items()[0].(models.Item[testNamed]).Title(); got != "Fight" {
		t.Fatalf("got title %q", got)
	}
}
