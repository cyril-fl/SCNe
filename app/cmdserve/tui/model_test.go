package cmdservetui

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"

	"solopg/app/shared/components/models"
	"solopg/app/shared/messages"
)

type testSubmodel struct {
	received *[]tea.Msg
}

func (testSubmodel) Init() tea.Cmd { return nil }
func (m testSubmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	*m.received = append(*m.received, msg)
	return m, nil
}
func (testSubmodel) View() tea.View { return tea.NewView("submodel") }

func newTestModel(received *[]tea.Msg) model {
	return NewModel(&models.ContextStepList[Context]{
		Steps: []models.Step[Context]{{Submodel: testSubmodel{received: received}}},
	})
}

func TestModelDelegatesToCurrentSubmodel(t *testing.T) {
	var received []tea.Msg
	m := newTestModel(&received)

	m.Update("event")

	if len(received) != 1 || received[0] != "event" {
		t.Fatalf("submodel received %v", received)
	}
	if got := m.View().Content; got != "submodel" {
		t.Fatalf("view = %q, want submodel view", got)
	}
}

func TestModelStoresResolvedErrorAndQuits(t *testing.T) {
	var received []tea.Msg
	m := newTestModel(&received)
	want := errors.New("stream failed")

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
	m := newTestModel(&received)

	_, quit := m.Update(tea.KeyPressMsg(tea.Key{Code: 'q', Mod: tea.ModCtrl}))

	if _, ok := quit().(tea.QuitMsg); !ok {
		t.Fatal("expected tea.Quit after ctrl+q")
	}
}
