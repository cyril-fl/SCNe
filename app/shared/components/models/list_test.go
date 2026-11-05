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

func TestListContentHeightFitsDescriptions(t *testing.T) {
	compact := NewList([]list.Item{NewItem("a", "", 1), NewItem("b", "", 2)}, 0, 0)
	if got := ListContentHeight(&compact); got != 2 {
		t.Fatalf("ListContentHeight() without descriptions = %d, want 2", got)
	}

	described := NewList([]list.Item{NewItem("a", "desc", 1), NewItem("b", "", 2)}, 0, 0)
	if got := ListContentHeight(&described); got != 6 {
		t.Fatalf("ListContentHeight() with a description = %d, want 6", got)
	}
}

func TestFitListCapsHeightAndFitsWidth(t *testing.T) {
	short := NewList([]list.Item{NewItem("a", "", 1), NewItem("longest", "", 2)}, 0, 0)
	FitList(&short, 5)
	if short.Height() != 2 || short.Width() != len("longest")+2 {
		t.Fatalf("FitList() short = %dx%d, want %dx2", short.Width(), short.Height(), len("longest")+2)
	}

	items := make([]list.Item, 0, 8)
	for i := range 8 {
		items = append(items, NewItem("x", "", i))
	}
	long := NewList(items, 0, 0)
	FitList(&long, 5)
	if long.Height() != 5 {
		t.Fatalf("FitList() long height = %d, want 5", long.Height())
	}
}
