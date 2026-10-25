package menuaction

import (
	"testing"

	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/form/field"
	"solopg/app/shared/components/models"

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
