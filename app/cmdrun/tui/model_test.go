package cmdruntui

import (
	"errors"
	"os"
	"testing"

	tea "charm.land/bubbletea/v2"

	"scne/app/shared/components/models"
	"scne/app/shared/messages"
	"scne/app/shared/services/logs"
)

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func TestMain(m *testing.M) {
	logs.Init(testLogRegistrar{})
	code := m.Run()
	logs.Init(nil)
	os.Exit(code)
}

type testSubmodel struct {
	received      *[]tea.Msg
	handlesEscape bool
}

func (testSubmodel) Init() tea.Cmd { return nil }
func (m testSubmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	*m.received = append(*m.received, msg)
	return m, nil
}
func (testSubmodel) View() tea.View        { return tea.NewView("submodel") }
func (m testSubmodel) HandlesEscape() bool { return m.handlesEscape }

func newTestModel(submodel testSubmodel) model {
	return NewModel(&models.ContextStepList[Context]{
		Steps: []models.Step[Context]{{Submodel: submodel}},
	})
}

func TestModelDelegatesToCurrentSubmodel(t *testing.T) {
	var received []tea.Msg
	m := newTestModel(testSubmodel{received: &received})

	m.Update("event")

	if len(received) != 1 || received[0] != "event" {
		t.Fatalf("submodel received %v", received)
	}
}

func TestModelStoresResolvedErrorAndQuits(t *testing.T) {
	var received []tea.Msg
	m := newTestModel(testSubmodel{received: &received})
	want := errors.New("step failed")

	_, quit := m.Update(messages.Resolved{Err: want})

	if !errors.Is(m.steps.GetErr(), want) {
		t.Fatalf("model error = %v, want %v", m.steps.GetErr(), want)
	}
	if _, ok := quit().(tea.QuitMsg); !ok {
		t.Fatal("expected tea.Quit after a resolved error")
	}
}

func TestModelQuitsOnCtrlQ(t *testing.T) {
	var received []tea.Msg
	m := newTestModel(testSubmodel{received: &received})

	_, quit := m.Update(tea.KeyPressMsg(tea.Key{Code: 'q', Mod: tea.ModCtrl}))

	if _, ok := quit().(tea.QuitMsg); !ok {
		t.Fatal("expected tea.Quit after ctrl+q")
	}
}

func TestModelEscapeQuitsUnlessSubmodelHandlesIt(t *testing.T) {
	esc := tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape})

	var received []tea.Msg
	m := newTestModel(testSubmodel{received: &received})
	if _, quit := m.Update(esc); quit == nil {
		t.Fatal("expected tea.Quit when the submodel does not handle escape")
	} else if _, ok := quit().(tea.QuitMsg); !ok {
		t.Fatal("expected tea.Quit when the submodel does not handle escape")
	}

	received = nil
	m = newTestModel(testSubmodel{received: &received, handlesEscape: true})
	m.Update(esc)
	if len(received) != 1 {
		t.Fatal("escape was not forwarded to the submodel that handles it")
	}
}
