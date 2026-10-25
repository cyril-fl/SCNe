package models

import (
	"errors"
	"os"
	"testing"

	"solopg/app/shared/messages"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/types/step"

	tea "charm.land/bubbletea/v2"
)

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func TestMain(m *testing.M) {
	logs.Init(testLogRegistrar{})
	code := m.Run()
	logs.Init(nil)
	os.Exit(code)
}

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
	m.ResolveSubmodel(messages.Resolved{Completed: true, Value: "ok"})
	if !resolved || m.CurrentIndex != 1 {
		t.Fatal("resolution did not resolve and advance to the next step")
	}
}

func TestStepModelResolutionErrorQuits(t *testing.T) {
	m := NewStepModel(&step.List[string]{Steps: []step.Step[string]{{}}})
	cmd := m.ResolveSubmodel(messages.Resolved{Err: errors.New("failed"), Completed: true})
	if cmd == nil || !m.HasErr() {
		t.Fatal("resolution error was not stored")
	}
}

func TestBind(t *testing.T) {
	s := Bind(testStepModel{}, "name", func(ctx *string, v string) { *ctx = v })
	ctx := ""
	if err := s.Resolve(&ctx, "ok"); err != nil || ctx != "ok" {
		t.Fatal("bound value was not applied to the context")
	}
	if err := s.Resolve(&ctx, 42); err == nil || ctx != "ok" {
		t.Fatal("wrong value type should fail without touching the context")
	}
}
