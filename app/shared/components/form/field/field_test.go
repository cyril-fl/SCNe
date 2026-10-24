package field

import (
	"errors"
	"strings"
	"testing"

	"solopg/app/shared/components/models"
	"solopg/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)

var errInvalid = errors.New("invalid value")

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

// - Field - //

func TestFieldErrorsAccumulateAndReset(t *testing.T) {
	f := newField("name", "", INPUT, "", nil, false, false)
	first, second := errors.New("first"), errors.New("second")

	f.SetError(first)
	f.SetError(second)
	if err := f.GetError(); !errors.Is(err, first) || !errors.Is(err, second) {
		t.Fatalf("GetError() = %v, want both errors", err)
	}

	f.ResetError()
	if err := f.GetError(); err != nil {
		t.Fatalf("GetError() after reset = %v, want nil", err)
	}
}

// - Text - //
func TestTextFieldUsesDefaultValue(t *testing.T) {
	f := TextField(TextTemplate[string]{Template: Template[string]{ID: "name", Defaultvalue: "Ada"}})

	if got := f.GetValueAsString(); got != "Ada" {
		t.Fatalf("GetValueAsString() = %q, want Ada", got)
	}
}

func TestTextFieldIntDefaultIsRendered(t *testing.T) {
	f := TextField(TextTemplate[int]{Template: Template[int]{ID: "age", Defaultvalue: 42}})

	if got := f.GetValueAsString(); got != "42" {
		t.Fatalf("GetValueAsString() = %q, want 42", got)
	}
}

func TestTextFieldFocusAndBlur(t *testing.T) {
	f := TextField(TextTemplate[string]{Template: Template[string]{ID: "name"}})

	f.Focus()
	if !f.IsFocused() {
		t.Fatal("Focus() did not focus field")
	}
	f.Blur()
	if f.IsFocused() {
		t.Fatal("Blur() did not blur field")
	}
}

func TestTextFieldTypingUpdatesValue(t *testing.T) {
	f := TextField(TextTemplate[string]{Template: Template[string]{ID: "name"}})
	f.Focus()

	f.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})

	if got := f.GetValueAsString(); got != "a" {
		t.Fatalf("value after typing = %q, want a", got)
	}
}

func TestTextFieldNavigationKeysReturnCommands(t *testing.T) {
	f := TextField(TextTemplate[string]{Template: Template[string]{ID: "name"}})

	for _, code := range []rune{tea.KeyEnter, tea.KeyUp, tea.KeyDown} {
		if _, cmd := f.Update(tea.KeyPressMsg{Code: code}); cmd == nil {
			t.Fatalf("Update(%q) returned nil command", code)
		}
	}
}

func TestTextFieldValidateRequired(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })

	f := TextField(TextTemplate[string]{Template: Template[string]{ID: "name", Defaultvalue: "  ", Required: true}})

	if err := f.Validate(); err == nil {
		t.Fatal("Validate() accepted blank required value")
	}
}

func TestTextFieldValidateRecordsValidatorError(t *testing.T) {
	f := TextField(TextTemplate[string]{Template: Template[string]{
		ID:           "name",
		Defaultvalue: "bad",
		Validator:    func(string) error { return errInvalid },
	}})

	if err := f.Validate(); !errors.Is(err, errInvalid) {
		t.Fatalf("Validate() = %v, want %v", err, errInvalid)
	}
	if !errors.Is(f.GetError(), errInvalid) {
		t.Fatalf("GetError() = %v, want %v", f.GetError(), errInvalid)
	}
	if !strings.Contains(f.View().Content, errInvalid.Error()) {
		t.Fatal("View() does not show validation error")
	}
}

func TestTextFieldResetRestoresDefault(t *testing.T) {
	f := TextField(TextTemplate[string]{Template: Template[string]{ID: "name", Defaultvalue: "Ada"}})
	f.Focus()
	f.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	f.SetError(errInvalid)

	f.Reset()

	if got := f.GetValueAsString(); got != "Ada" {
		t.Fatalf("value after reset = %q, want Ada", got)
	}
	if f.GetError() != nil || f.IsFocused() {
		t.Fatal("Reset() did not clear errors and focus")
	}
}

func TestTextFieldViewLabel(t *testing.T) {
	labelled := TextField(TextTemplate[string]{Template: Template[string]{ID: "name", Label: "unknownlabel"}})
	unlabelled := TextField(TextTemplate[string]{Template: Template[string]{ID: "name", Label: "unknownlabel", NoLabel: true}})

	if !strings.Contains(labelled.View().Content, "unknownlabel") {
		t.Fatal("View() does not show label")
	}
	if strings.Contains(unlabelled.View().Content, "unknownlabel") {
		t.Fatal("View() shows label with NoLabel")
	}
}

// - Select - //
func testSelect(template Template[string]) *selectField[string] {
	return SelectField(SelectTemplate[string]{
		Template: template,
		Options: []models.Item[string]{
			models.NewItem("First", "", "first"),
			models.NewItem("Second", "", "second"),
		},
	})
}

func TestSelectFieldSelectsDefaultValue(t *testing.T) {
	f := testSelect(Template[string]{ID: "choice", Defaultvalue: "second"})

	if got := f.GetValueAsString(); got != "second" {
		t.Fatalf("GetValueAsString() = %q, want second", got)
	}
}

func TestSelectFieldResetSelectsFirstOption(t *testing.T) {
	f := testSelect(Template[string]{ID: "choice", Defaultvalue: "second"})
	f.Focus()
	f.SetError(errInvalid)

	f.Reset()

	if got := f.GetValueAsString(); got != "first" {
		t.Fatalf("value after reset = %q, want first", got)
	}
	if f.GetError() != nil || f.IsFocused() {
		t.Fatal("Reset() did not clear errors and focus")
	}
}

func TestSelectFieldValidateRecordsValidatorError(t *testing.T) {
	var validated string
	f := testSelect(Template[string]{
		ID:           "choice",
		Defaultvalue: "second",
		Validator: func(value string) error {
			validated = value
			return errInvalid
		},
	})

	if err := f.Validate(); !errors.Is(err, errInvalid) {
		t.Fatalf("Validate() = %v, want %v", err, errInvalid)
	}
	if validated != "second" {
		t.Fatalf("validator received %q, want second", validated)
	}
	if !strings.Contains(f.View().Content, errInvalid.Error()) {
		t.Fatal("View() does not show validation error")
	}
}

func TestSelectFieldEnterReturnsCommand(t *testing.T) {
	f := testSelect(Template[string]{ID: "choice"})

	if _, cmd := f.Update(tea.KeyPressMsg{Code: tea.KeyEnter}); cmd == nil {
		t.Fatal("Update(enter) returned nil command")
	}
}

// - Textarea - //
func TestTextareaFieldUsesDefaultValue(t *testing.T) {
	f := TextareaField(TextareaTemplate{Template: Template[string]{ID: "message", Defaultvalue: "Ada"}})

	if got := f.GetValueAsString(); got != "Ada" {
		t.Fatalf("GetValueAsString() = %q, want Ada", got)
	}
}

func TestTextareaFieldFocusAndBlur(t *testing.T) {
	f := TextareaField(TextareaTemplate{Template: Template[string]{ID: "message"}})

	f.Focus()
	if !f.IsFocused() {
		t.Fatal("Focus() did not focus field")
	}
	f.Blur()
	if f.IsFocused() {
		t.Fatal("Blur() did not blur field")
	}
}

func TestTextareaFieldTypingUpdatesValue(t *testing.T) {
	f := TextareaField(TextareaTemplate{Template: Template[string]{ID: "message"}})
	f.Focus()

	f.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})

	if got := f.GetValueAsString(); got != "a" {
		t.Fatalf("value after typing = %q, want a", got)
	}
}

func TestTextareaFieldEnterReturnsCommandWithoutNewline(t *testing.T) {
	f := TextareaField(TextareaTemplate{Template: Template[string]{ID: "message", Defaultvalue: "hi"}})
	f.Focus()

	if _, cmd := f.Update(tea.KeyPressMsg{Code: tea.KeyEnter}); cmd == nil {
		t.Fatal("Update(enter) returned nil command")
	}
	if got := f.GetValueAsString(); got != "hi" {
		t.Fatalf("value after enter = %q, want hi", got)
	}
}

func TestTextareaFieldArrowsStayInInput(t *testing.T) {
	f := TextareaField(TextareaTemplate{Template: Template[string]{ID: "message", Defaultvalue: "hi"}})
	f.Focus()

	for _, code := range []rune{tea.KeyUp, tea.KeyDown} {
		if _, cmd := f.Update(tea.KeyPressMsg{Code: code}); cmd != nil {
			t.Fatalf("Update(%q) returned a form command", code)
		}
	}
	if got := f.GetValueAsString(); got != "hi" {
		t.Fatalf("value after arrows = %q, want hi", got)
	}
}

func TestTextareaFieldAltEnterIgnoredWhenEmpty(t *testing.T) {
	f := TextareaField(TextareaTemplate{Template: Template[string]{ID: "message"}})
	f.Focus()

	f.Update(tea.KeyPressMsg{Code: tea.KeyEnter, Mod: tea.ModAlt})

	if got := f.GetValueAsString(); got != "" {
		t.Fatalf("value after alt+enter = %q, want empty", got)
	}
}

func TestTextareaFieldAltEnterInsertsNewline(t *testing.T) {
	f := TextareaField(TextareaTemplate{Template: Template[string]{ID: "message"}})
	f.Focus()

	f.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	f.Update(tea.KeyPressMsg{Code: tea.KeyEnter, Mod: tea.ModAlt})
	f.Update(tea.KeyPressMsg{Code: 'b', Text: "b"})

	if got := f.GetValueAsString(); got != "a\nb" {
		t.Fatalf("value after alt+enter = %q, want %q", got, "a\nb")
	}
}

func TestTextareaFieldValidateRequired(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })

	f := TextareaField(TextareaTemplate{Template: Template[string]{ID: "message", Defaultvalue: "  ", Required: true}})

	if err := f.Validate(); err == nil {
		t.Fatal("Validate() accepted blank required value")
	}
}

func TestTextareaFieldValidateRecordsValidatorError(t *testing.T) {
	f := TextareaField(TextareaTemplate{Template: Template[string]{
		ID:           "message",
		Defaultvalue: "bad",
		Validator:    func(string) error { return errInvalid },
	}})

	if err := f.Validate(); !errors.Is(err, errInvalid) {
		t.Fatalf("Validate() = %v, want %v", err, errInvalid)
	}
	if !strings.Contains(f.View().Content, errInvalid.Error()) {
		t.Fatal("View() does not show validation error")
	}
}

func TestTextareaFieldResetRestoresDefault(t *testing.T) {
	f := TextareaField(TextareaTemplate{Template: Template[string]{ID: "message", Defaultvalue: "Ada"}})
	f.Focus()
	f.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	f.SetError(errInvalid)

	f.Reset()

	if got := f.GetValueAsString(); got != "Ada" {
		t.Fatalf("value after reset = %q, want Ada", got)
	}
	if f.GetError() != nil || f.IsFocused() {
		t.Fatal("Reset() did not clear errors and focus")
	}
}

func TestTextareaFieldViewLabel(t *testing.T) {
	labelled := TextareaField(TextareaTemplate{Template: Template[string]{ID: "message", Label: "unknownlabel"}})
	unlabelled := TextareaField(TextareaTemplate{Template: Template[string]{ID: "message", Label: "unknownlabel", NoLabel: true}})

	if !strings.Contains(labelled.View().Content, "unknownlabel") {
		t.Fatal("View() does not show label")
	}
	if strings.Contains(unlabelled.View().Content, "unknownlabel") {
		t.Fatal("View() shows label with NoLabel")
	}
}
