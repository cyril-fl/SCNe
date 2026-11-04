package menuaction

import (
	"strings"
	"testing"

	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/form"
	"scne/app/shared/components/form/field"
	"scne/app/shared/components/models"

	tea "charm.land/bubbletea/v2"
)

func newTestPage(submitted *string) *actionPage {
	return newPage(PageTemplate{
		Label: "mock",
		NewForm: func() *form.Form {
			return form.New(
				field.SelectField(field.SelectTemplate[string]{Template: field.Template[string]{ID: "choice", Required: true}, Options: []models.Item[string]{models.NewItem("first", "", "first")}}),
			).Focus()
		},
		Submit: func(values *form.Form) tea.Cmd {
			*submitted = values.GetFieldValue("choice").(string)
			return nil
		},
	})
}

func TestPageBuildsFormOnEntry(t *testing.T) {
	var submitted string
	page := newTestPage(&submitted)

	if page.GetView() != "" {
		t.Fatal("expected no form before the page is entered")
	}

	page.SetSize(size.New(size.Template{Width: 40, Height: 20}))
	if page.GetView() == "" {
		t.Fatal("expected the form after the page is entered")
	}
}

func TestPageValidateSubmitsAndCloses(t *testing.T) {
	var submitted string
	page := newTestPage(&submitted)
	page.SetSize(size.New(size.Template{Width: 40, Height: 20}))

	cmd := page.Update(form.Validate{})

	if submitted != "first" {
		t.Fatalf("submitted = %q, want first", submitted)
	}
	if cmd == nil || page.GetView() != "" {
		t.Fatal("expected a close command and the form to be dropped")
	}
}

func TestPageReopensWithCurrentOptions(t *testing.T) {
	option := "first"
	page := newPage(PageTemplate{
		NewForm: func() *form.Form {
			return form.New(
				field.SelectField(field.SelectTemplate[string]{Template: field.Template[string]{ID: "choice", Required: true}, Options: []models.Item[string]{models.NewItem(option, "", option)}}),
			).Focus()
		},
	})
	pageSize := size.New(size.Template{Width: 40, Height: 20})
	page.SetSize(pageSize)

	closing := page.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	// The layout resizes the current view on each render, before form.Close is handled.
	page.SetSize(pageSize)
	page.Update(closing())

	option = "second"
	page.SetSize(pageSize)

	if view := page.GetView(); !strings.Contains(view, "Second") {
		t.Fatalf("reopened form does not list the new option:\n%s", view)
	}
}

func TestPageEscapeCloses(t *testing.T) {
	var submitted string
	page := newTestPage(&submitted)
	page.SetSize(size.New(size.Template{Width: 40, Height: 20}))

	cmd := page.Update(tea.KeyPressMsg{Code: tea.KeyEscape})

	if cmd == nil {
		t.Fatal("expected a command")
	}
	if _, ok := cmd().(form.Close); !ok {
		t.Fatalf("expected form.Close, got %#v", cmd())
	}
	if submitted != "" {
		t.Fatal("escape must not submit")
	}
}
