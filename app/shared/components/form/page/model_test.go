package formpage

import (
	"errors"
	"testing"

	"scne/app/cmdrun/types/id"
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/form"
	"scne/app/shared/components/form/field"

	tea "charm.land/bubbletea/v2"
)

func testForm(values map[string]string) *form.Form {
	return form.New(field.TextField(field.TextTemplate[string]{
		ID: "name", Defaultvalue: values["name"],
	})).Focus()
}

func testKey(code rune, mod tea.KeyMod) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code, Mod: mod}
}

func TestPageSetSizeResizesList(t *testing.T) {
	page := New(Template{
		Items:   func() []Item { return nil },
		NewForm: testForm,
		Save:    func(id.ID, map[string]string) error { return nil },
	})

	page.SetSize(size.New(size.Template{Width: 42, Height: 20}))
	if got := page.GetList().Width(); got != 42 {
		t.Fatalf("expected list width 42, got %d", got)
	}
}

func TestPageCtrlNOpensNewFormFromEmptyPage(t *testing.T) {
	page := New(Template{
		Items:   func() []Item { return nil },
		NewForm: testForm,
		Save:    func(id.ID, map[string]string) error { return nil },
	})

	page.Update(testKey('n', tea.ModCtrl))
	if !page.IsFormOpen() {
		t.Fatal("expected Ctrl+N to open a form")
	}
	if page.GetForm().GetFieldValue("name") != "" {
		t.Fatal("expected a new form to have empty values")
	}
}

func TestPageCtrlEOpensSelectedEntry(t *testing.T) {
	page := New(Template{
		Items:   func() []Item { return []Item{{Name: "first", Values: map[string]string{"name": "loaded"}}} },
		NewForm: testForm,
		Save:    func(id.ID, map[string]string) error { return nil },
	})

	page.Update(testKey('e', tea.ModCtrl))
	if got := page.GetForm().GetFieldValue("name"); got != "loaded" {
		t.Fatalf("expected selected value to be loaded, got %v", got)
	}
}

func TestPageCtrlEDoesNotOpenNewFormWhenListIsEmpty(t *testing.T) {
	page := New(Template{
		Items:   func() []Item { return nil },
		NewForm: testForm,
		Save:    func(id.ID, map[string]string) error { return nil },
	})

	page.Update(testKey('e', tea.ModCtrl))
	if page.IsFormOpen() {
		t.Fatal("expected Ctrl+E on an empty list not to open a form")
	}
}

func TestPageCtrlEDoesNotOpenLockedEntry(t *testing.T) {
	page := New(Template{
		Items: func() []Item {
			return []Item{{Name: "template", Locked: true, Values: map[string]string{"name": "template"}}}
		},
		NewForm: testForm,
		Save:    func(id.ID, map[string]string) error { return nil },
	})

	if page.EditItem() {
		t.Fatal("expected a locked entry not to open for editing")
	}
	if page.IsFormOpen() {
		t.Fatal("expected a locked entry not to open a form")
	}
}

func TestPageCtrlNOpensNewForm(t *testing.T) {
	page := New(Template{
		Items:   func() []Item { return []Item{{Name: "first"}} },
		NewForm: testForm,
		Save:    func(id.ID, map[string]string) error { return nil },
	})

	page.Update(testKey('n', tea.ModCtrl))
	if !page.IsFormOpen() {
		t.Fatal("expected Ctrl+N to open a new form")
	}
}

func TestPageListNavigationAndEscape(t *testing.T) {
	page := New(Template{
		Items:   func() []Item { return []Item{{Name: "first"}, {Name: "second"}} },
		NewForm: testForm,
		Save:    func(id.ID, map[string]string) error { return nil },
	})

	page.Update(testKey(tea.KeyDown, 0))
	if page.GetList().Index() != 1 {
		t.Fatalf("expected list cursor 1, got %d", page.GetList().Index())
	}
	page.Update(testKey('e', tea.ModCtrl))
	page.Update(testKey(tea.KeyEsc, 0))
	if page.IsFormOpen() {
		t.Fatal("expected Escape to close the form")
	}
	if page.HandleEscape() {
		t.Fatal("expected second Escape to be handled by the owning menu")
	}
}

func TestPageSubmitCallsSave(t *testing.T) {
	called := false
	page := New(Template{
		Items:   func() []Item { return nil },
		NewForm: testForm,
		Save: func(entryID id.ID, values map[string]string) error {
			called = entryID.IsZero() && values["name"] == ""
			return nil
		},
	})

	page.EditNewItem()
	page.Update(form.Validate{})
	if !called {
		t.Fatal("expected form submission callback")
	}
	if page.IsFormOpen() {
		t.Fatal("expected form to close after saving")
	}
}

func TestPageCtrlEnterValidates(t *testing.T) {
	called := false
	page := New(Template{
		Items:   func() []Item { return nil },
		NewForm: testForm,
		Save: func(entryID id.ID, values map[string]string) error {
			called = true
			return nil
		},
	})

	page.EditNewItem()
	cmd := page.Update(testKey(tea.KeyEnter, tea.ModCtrl))
	if cmd == nil {
		t.Fatal("expected Ctrl+Enter to request validation")
	}
	page.Update(cmd())
	if !called {
		t.Fatal("expected Ctrl+Enter to submit the form")
	}
}

func TestPageLastFieldValidates(t *testing.T) {
	called := false
	page := New(Template{
		Items:   func() []Item { return nil },
		NewForm: testForm,
		Save: func(entryID id.ID, values map[string]string) error {
			called = true
			return nil
		},
	})

	page.EditNewItem()
	cmd := page.Update(form.NextField{})
	if cmd == nil {
		t.Fatal("expected last field to request validation")
	}
	page.Update(cmd())
	if !called {
		t.Fatal("expected last field transition to submit the form")
	}
}

func TestPageSaveReceivesSelectedEntryIDWhenFiltered(t *testing.T) {
	alpha, beta := id.New(), id.New()
	var saved id.ID
	page := New(Template{
		Items: func() []Item {
			return []Item{{ID: alpha, Name: "alpha"}, {ID: beta, Name: "beta"}}
		},
		NewForm: testForm,
		Save: func(entryID id.ID, values map[string]string) error {
			saved = entryID
			return nil
		},
	})

	page.GetList().SetFilterText("beta")
	if !page.EditItem() {
		t.Fatal("expected filtered entry to open for editing")
	}
	page.Update(form.Validate{})

	if saved != beta {
		t.Fatalf("Save received %v, want %v", saved, beta)
	}
}

func TestPageSaveErrorKeepsFormOpen(t *testing.T) {
	page := New(Template{
		Items:   func() []Item { return nil },
		NewForm: testForm,
		Save:    func(id.ID, map[string]string) error { return errors.New("save failed") },
	})

	page.EditNewItem()
	if cmd := page.Update(form.Validate{}); cmd == nil {
		t.Fatal("expected a failure command when Save fails")
	}
	if !page.IsFormOpen() || !page.GetForm().HasErr() {
		t.Fatal("expected form to stay open with the save error")
	}
}

func TestPageCtrlEnterRunsActionOnSelectedEntry(t *testing.T) {
	entryID := id.New()
	var got id.ID
	page := New(Template{
		Items:   func() []Item { return []Item{{ID: entryID, Name: "first"}} },
		NewForm: testForm,
		Save:    func(id.ID, map[string]string) error { return nil },
		Action:  &Action{Run: func(selected id.ID) { got = selected }},
	})

	page.Update(testKey(tea.KeyEnter, tea.ModCtrl))
	if got != entryID {
		t.Fatalf("expected action to receive %v, got %v", entryID, got)
	}
	if page.IsFormOpen() {
		t.Fatal("expected Ctrl+Enter not to open a form")
	}
}
