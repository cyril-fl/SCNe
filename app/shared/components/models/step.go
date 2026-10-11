package models

import (
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/types/primitive"
	"solopg/app/shared/types/step"

	tea "charm.land/bubbletea/v2"
)

type ContextStepList[T any] = step.List[T]
type Step[T any] = step.Step[T]

type StepModel[T any] struct {
	primitive.Fallible
	*ContextStepList[T]
	Context T
	size    *tea.WindowSizeMsg
}

func NewStepModel[T any](steps *ContextStepList[T]) StepModel[T] {
	return StepModel[T]{
		ContextStepList: steps,
	}
}

// Getters & Setters
func (m *StepModel[T]) SetSize(size *tea.WindowSizeMsg) {
	if size == nil {
		return
	}
	m.size = size
}

func (m *StepModel[T]) GetSize() *tea.WindowSizeMsg {
	return m.size
}

// Methods
func (m *StepModel[T]) ResolveSubmodel(resolution sharedtui.ResolutionMsg) tea.Cmd {
	if resolution.Err != nil {
		m.SetErr(resolution.Err)
		return tea.Quit
	}

	if !resolution.Completed {
		return nil
	}

	if err := m.handleResolution(resolution.Value); err != nil {
		m.SetErr(err)
		return tea.Quit
	}

	return m.ForwardNext()
}

func (m *StepModel[T]) UpdateSubmodel(msg tea.Msg) tea.Cmd {
	current := m.GetCurrentSubmodel()
	if current == nil {
		return nil
	}

	var cmd tea.Cmd
	current, cmd = current.Update(msg)

	m.SetCurrentSubmodel(current)

	return cmd
}

func (m *StepModel[T]) SelectNext() bool {
	for {
		next := m.NextStep()
		if next == nil {
			return false
		}

		if next.Skip == nil || !next.Skip(&m.Context) {
			return true
		}
	}
}

func (m *StepModel[T]) ForwardNext() tea.Cmd {
	if !m.SelectNext() {
		return tea.Quit
	}

	submodel := m.GetCurrentSubmodel()
	if submodel == nil {
		return tea.Quit
	}

	initCmd := submodel.Init()
	if m.size == nil {
		return initCmd
	}

	return tea.Sequence(
		initCmd,
		sharedtui.SendSizeMsg(m.size),
	)
}

func (m *StepModel[T]) handleResolution(value any) error {
	if current := m.CurrentStep(); current != nil && current.Resolve != nil {
		return current.Resolve(&m.Context, value)
	}

	return nil
}

// Helpers
