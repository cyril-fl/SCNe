package navigation

import (
	"strings"
	"testing"

	"solopg/app/cmdrun/types/size"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type testItem struct {
	content string
	list    list.Model
}

func newTestItem(content string) *testItem {
	return &testItem{
		content: content,
		list:    list.New(nil, list.NewDefaultDelegate(), 0, 0),
	}
}

func (i *testItem) ID() string { return i.content }

func (i *testItem) SetFocus(bool) {}

func (i *testItem) GetList() *list.Model { return &i.list }

func (i *testItem) SetList(value *list.Model) { i.list = *value }

func (i *testItem) SetSize(size.Size) {}

func (i *testItem) GetMenuView() string { return i.content }

func (i *testItem) GetView() string { return i.content }

func (i *testItem) GetFooter() []string { return nil }

func (i *testItem) HandleUpdate(tea.Msg) tea.Cmd { return nil }

func newTestNavigation(contents []string, viewportHeight int) *Navigation {
	parts := make([]Item, 0, len(contents))
	for _, content := range contents {
		parts = append(parts, newTestItem(content))
	}

	navigation := New(Template{
		Static: []string{"static"},
		Parts:  parts,
	})
	navigation.SetSize(size.New(size.Template{
		Width:  24,
		Height: viewportHeight + 2,
	}))
	navigation.view.SetHeight(viewportHeight)

	return navigation
}

func TestGetViewKeepsStaticContentOutsideViewport(t *testing.T) {
	navigation := newTestNavigation([]string{"one", "two", "three", "four", "five"}, 2)
	navigation.GetView()

	if strings.Contains(navigation.view.GetContent(), "static") {
		t.Fatal("expected static content to be excluded from viewport content")
	}
}
