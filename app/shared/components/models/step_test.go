package models

import (
	"errors"
	"testing"

	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/types/step"

	tea "charm.land/bubbletea/v2"
)

type testStepModel struct{}

func (testStepModel) Init() tea.Cmd                         { return nil }
func (m testStepModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (testStepModel) View() tea.View                        { return tea.NewView("") }

func TestStepModelResolveAndForward(t *testing.T) {
	resolved := false
	steps := &step.List[string]{Steps: []step.Step[string]{
		{Resolve: func(ctx *string, value any) error { resolved = value == "ok"; return nil }},
		{Submodel: testStepModel{}},
	}}
	m := NewStepModel(steps)
	m.ResolveSubmodel(sharedtui.ResolutionMsg{Completed: true, Value: "ok"})
	if !resolved || m.CurrentIndex != 1 {
		t.Fatal("resolution did not resolve and advance to the next step")
	}
}

func TestStepModelResolutionErrorQuits(t *testing.T) {
	m := NewStepModel(&step.List[string]{Steps: []step.Step[string]{{}}})
	cmd := m.ResolveSubmodel(sharedtui.ResolutionMsg{Err: errors.New("failed"), Completed: true})
	if cmd == nil || !m.HasErr() {
		t.Fatal("resolution error was not stored")
	}
}
