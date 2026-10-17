package formpage

import (
	"testing"

	"solopg/app/shared/components/form"
	"solopg/app/shared/components/form/field"

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

func TestPageCtrlNOpensNewFormFromEmptyPage(t *testing.T) {
	page := New(Template{
		Items:   func() []Item { return nil },
		NewForm: testForm,
		Save:    func(int, map[string]string) error { return nil },
	})

	page.Update(testKey('n', tea.ModCtrl))
	if !page.IsFormOpen() {
		t.Fatal("expected Ctrl+N to open a form")
	}
	if page.GetForm().GetFieldValue("name") != "" {
		t.Fatal("expected a new form to have empty values")
	}
}

func TestPageEnterOpensSelectedEntry(t *testing.T) {
	page := New(Template{
		Items:   func() []Item { return []Item{{Name: "first", Values: map[string]string{"name": "loaded"}}} },
		NewForm: testForm,
		Save:    func(int, map[string]string) error { return nil },
	})

	page.Update(testKey(tea.KeyEnter, tea.ModAlt))
	if got := page.GetForm().GetFieldValue("name"); got != "loaded" {
		t.Fatalf("expected selected value to be loaded, got %v", got)
	}
}

func TestPageAltEnterDoesNotOpenNewFormWhenListIsEmpty(t *testing.T) {
	page := New(Template{
		Items:   func() []Item { return nil },
		NewForm: testForm,
		Save:    func(int, map[string]string) error { return nil },
	})

	page.Update(testKey(tea.KeyEnter, tea.ModAlt))
	if page.IsFormOpen() {
		t.Fatal("expected Alt+Enter on an empty list not to open a form")
	}
}

func TestPageAltEnterDoesNotOpenLockedEntry(t *testing.T) {
	page := New(Template{
		Items: func() []Item {
			return []Item{{Name: "template", Locked: true, Values: map[string]string{"name": "template"}}}
		},
		NewForm: testForm,
		Save:    func(int, map[string]string) error { return nil },
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
		Save:    func(int, map[string]string) error { return nil },
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
		Save:    func(int, map[string]string) error { return nil },
	})

	page.Update(testKey(tea.KeyDown, 0))
	if page.GetList().Index() != 1 {
		t.Fatalf("expected list cursor 1, got %d", page.GetList().Index())
	}
	page.Update(testKey(tea.KeyEnter, tea.ModAlt))
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
		Save: func(index int, values map[string]string) error {
			called = index == -1 && values["name"] == ""
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

func TestPageAltEnterValidates(t *testing.T) {
	called := false
	page := New(Template{
		Items:   func() []Item { return nil },
		NewForm: testForm,
		Save: func(index int, values map[string]string) error {
			called = true
			return nil
		},
	})

	page.EditNewItem()
	cmd := page.Update(testKey(tea.KeyEnter, tea.ModAlt))
	if cmd == nil {
		t.Fatal("expected Alt+Enter to request validation")
	}
	page.Update(cmd())
	if !called {
		t.Fatal("expected Alt+Enter to submit the form")
	}
}

func TestPageLastFieldValidates(t *testing.T) {
	called := false
	page := New(Template{
		Items:   func() []Item { return nil },
		NewForm: testForm,
		Save: func(index int, values map[string]string) error {
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
