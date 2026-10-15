package models

import (
	"testing"

	"charm.land/bubbles/v2/list"
)

func TestItemExposesListFields(t *testing.T) {
	item := NewItem("Title", "Description", 42)
	if item.Title() != "Title" || item.Description() != "Description" || item.FilterValue() != "Title" || item.Value() != 42 {
		t.Fatalf("unexpected item: %#v", item)
	}
}

func TestConfigureListHidesOptionalElements(t *testing.T) {
	menu := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	ConfigureList(&menu)
	if menu.ShowFilter() || menu.ShowPagination() || menu.ShowHelp() || menu.ShowStatusBar() || menu.ShowTitle() {
		t.Fatal("list optional elements were not hidden")
	}
}
