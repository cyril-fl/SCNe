package primitive

import (
	"testing"

	"charm.land/bubbles/v2/list"

	"solopg/app/cmdrun/types/size"
)

func TestListedExposesTheProvidedList(t *testing.T) {
	model := list.New(nil, list.NewDefaultDelegate(), 12, 8)
	listed := Listed(model)

	if got := listed.GetList(); got.Width() != 12 || got.Height() != 8 {
		t.Fatalf("listed dimensions = %dx%d, want 12x8", got.Width(), got.Height())
	}
}

func TestListedGetListAllowsMutation(t *testing.T) {
	listed := Listed(list.New(nil, list.NewDefaultDelegate(), 0, 0))

	listed.GetList().SetSize(20, 10)

	if got := listed.GetList(); got.Width() != 20 || got.Height() != 10 {
		t.Fatalf("list dimensions after mutation = %dx%d, want 20x10", got.Width(), got.Height())
	}
}

func TestListedSetListReplacesTheList(t *testing.T) {
	listed := Listed(list.New(nil, list.NewDefaultDelegate(), 1, 2))
	replacement := list.New(nil, list.NewDefaultDelegate(), 30, 15)

	listed.SetList(&replacement)

	if got := listed.GetList(); got.Width() != 30 || got.Height() != 15 {
		t.Fatalf("replaced list dimensions = %dx%d, want 30x15", got.Width(), got.Height())
	}
}

func TestListedSetListSizeUsesTemplateDimensions(t *testing.T) {
	listed := Listed(list.New(nil, list.NewDefaultDelegate(), 0, 0))

	listed.SetListSize(size.Template{Width: 42, Height: 20})

	if got := listed.GetList(); got.Width() != 42 || got.Height() != 20 {
		t.Fatalf("list dimensions = %dx%d, want 42x20", got.Width(), got.Height())
	}
}

func TestListedSetListSizeClampsNegativeDimensions(t *testing.T) {
	listed := Listed(list.New(nil, list.NewDefaultDelegate(), 10, 10))

	listed.SetListSize(size.Template{Width: -5, Height: -3})

	if got := listed.GetList(); got.Width() != 0 || got.Height() != 0 {
		t.Fatalf("list dimensions = %dx%d, want 0x0", got.Width(), got.Height())
	}
}
