package form

import (
	"errors"
	"strings"
	"testing"

	"solopg/app/shared/services/logs"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

var errInvalid = errors.New("invalid value")

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

// testField is a minimal Field: the field package imports form, so it cannot be used here.
type testField struct {
	id          string
	value       string
	focus       bool
	err         []error
	validateErr error
	lastMsg     tea.Msg
}

func (f *testField) ID() string               { return f.id }
func (f *testField) Label() string            { return f.id }
func (f *testField) Focus() tea.Cmd           { f.focus = true; return nil }
func (f *testField) Blur()                    { f.focus = false }
func (f *testField) IsFocused() bool          { return f.focus }
func (f *testField) Value() any               { return f.value }
func (f *testField) GetValueAsString() string { return f.value }
func (f *testField) Reset()                   { f.value = "" }
func (f *testField) Validate() error          { return f.validateErr }
func (f *testField) SetError(err error)       { f.err = append(f.err, err) }
func (f *testField) GetError() error          { return errors.Join(f.err...) }
func (f *testField) ResetError()              { f.err = nil }
func (f *testField) Init() tea.Cmd            { return nil }
func (f *testField) View() tea.View           { return tea.NewView("view:" + f.id) }
func (f *testField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	f.lastMsg = msg
	return f, nil
}

func testForm() (*Form, *testField, *testField) {
	first, second := &testField{id: "first", value: "a"}, &testField{id: "second", value: "b"}
	return New(first, second), first, second
}

func TestFormValues(t *testing.T) {
	f, _, _ := testForm()

	if got := f.GetFieldValue("second"); got != "b" {
		t.Fatalf("GetFieldValue(second) = %v, want b", got)
	}
	if got := f.GetFieldValue("unknown"); got != nil {
		t.Fatalf("GetFieldValue(unknown) = %v, want nil", got)
	}
	if got := f.GetValuesAsMappedString(); got["first"] != "a" || got["second"] != "b" {
		t.Fatalf("GetValuesAsMappedString() = %v", got)
	}
}

func TestFormUpdateIgnoredWithoutFocus(t *testing.T) {
	f, first, second := testForm()

	if _, cmd := f.Update(NextField{}); cmd != nil {
		t.Fatal("Update() returned a command without a focused field")
	}
	if first.IsFocused() || second.IsFocused() {
		t.Fatal("Update() focused a field without a focused form")
	}
}

func TestFormNavigation(t *testing.T) {
	f, first, second := testForm()
	f.Focus()

	f.Update(NextField{})
	if !second.IsFocused() || first.IsFocused() {
		t.Fatal("NextField did not focus the second field")
	}
	f.Update(PreviousField{})
	if !first.IsFocused() || second.IsFocused() {
		t.Fatal("PreviousField did not focus the first field")
	}
	f.Update(PreviousField{})
	if !first.IsFocused() {
		t.Fatal("PreviousField on the first field did not stay on it")
	}
}

func TestFormLastFieldAutoSubmit(t *testing.T) {
	f, first, _ := testForm()
	f.focusOnIndex(1)

	_, cmd := f.Update(NextField{})
	if cmd == nil {
		t.Fatal("NextField on the last field did not request validation")
	}
	if _, ok := cmd().(Validate); !ok {
		t.Fatal("NextField on the last field did not send Validate")
	}

	f.DisableAutoSubmit()
	if _, cmd := f.Update(NextField{}); cmd != nil || !first.IsFocused() {
		t.Fatal("NextField without auto submit did not wrap to the first field")
	}
}

func TestFormKeyPressSendsMessages(t *testing.T) {
	f, first, _ := testForm()
	f.SetSendKey(plumtypes.KEY_ENTER).Focus()

	keys := map[tea.KeyPressMsg]tea.Msg{
		{Code: tea.KeyTab}:                    NextField{},
		{Code: tea.KeyTab, Mod: tea.ModShift}: PreviousField{},
		{Code: tea.KeyEnter}:                  Validate{},
	}
	for key, want := range keys {
		_, cmd := f.Update(key)
		if cmd == nil || cmd() != want {
			t.Fatalf("Update(%q) did not send %T", key.String(), want)
		}
	}

	key := tea.KeyPressMsg{Code: 'x', Text: "x"}
	f.Update(key)
	if first.lastMsg != key {
		t.Fatal("Update() did not forward other keys to the focused field")
	}
}

func TestFormValidateAndReset(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })

	f, first, second := testForm()
	second.validateErr = errInvalid
	f.focusOnIndex(1)

	f.Submit()
	if !f.HasErr() {
		t.Fatal("Submit() accepted an invalid field")
	}
	if !errors.Is(second.GetError(), errInvalid) {
		t.Fatalf("field error = %v, want %v", second.GetError(), errInvalid)
	}
	if !strings.Contains(f.View().Content, "view:first") {
		t.Fatal("View() does not render fields")
	}

	f.Reset()
	if f.HasErr() || second.GetError() != nil {
		t.Fatal("Reset() did not clear errors")
	}
	if first.value != "" || !first.IsFocused() || second.IsFocused() {
		t.Fatal("Reset() did not reset values and focus the first field")
	}
}

func TestFormGetFieldScrollArea(t *testing.T) {
	f, _, _ := testForm()
	vp := viewport.New(viewport.WithWidth(20), viewport.WithHeight(10))

	if top, height := f.GetFieldScrollArea(vp); top != -1 || height != 0 {
		t.Fatalf("GetFieldScrollArea() without focus = %d, %d; want -1, 0", top, height)
	}

	f.focusOnIndex(1)
	if top, height := f.GetFieldScrollArea(vp); top != 1 || height != 1 {
		t.Fatalf("GetFieldScrollArea() = %d, %d; want 1, 1", top, height)
	}
}
