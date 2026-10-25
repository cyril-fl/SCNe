package models

import (
	plumprimitive "github.com/cyril-fl/plum/primitive"
	"solopg/app/shared/messages"
	"solopg/app/shared/types/step"

	tea "charm.land/bubbletea/v2"
)

// ContextStepList is a step.List sharing context T.
type ContextStepList[T any] = step.List[T]

// Step is a step.Step sharing context T.
type Step[T any] = step.Step[T]

// StepModel runs a step list, each step resolving into the shared Context.
type StepModel[T any] struct {
	plumprimitive.Fallible
	*ContextStepList[T]
	Context T
	size    *tea.WindowSizeMsg
}

// NewStepModel returns a model running steps.
func NewStepModel[T any](steps *ContextStepList[T]) StepModel[T] {
	return StepModel[T]{
		ContextStepList: steps,
	}
}

// Getters & Setters

// SetSize stores size, ignoring nil, to resize each new step.
func (m *StepModel[T]) SetSize(size *tea.WindowSizeMsg) {
	if size == nil {
		return
	}
	m.size = size
}

// GetSize returns the last window size, or nil.
func (m *StepModel[T]) GetSize() *tea.WindowSizeMsg {
	return m.size
}

// Methods

// ResolveSubmodel applies a completed step to the context and moves to the next one.
// It quits on error, storing it.
func (m *StepModel[T]) ResolveSubmodel(resolution messages.Resolved) tea.Cmd {
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

// UpdateSubmodel forwards msg to the current step.
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

// SelectNext moves to the next step not skipped, false when none is left.
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

// ForwardNext starts the next step, resizing it, or quits when none is left.
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
		messages.Resize(m.size),
	)
}

func (m *StepModel[T]) handleResolution(value any) error {
	if current := m.CurrentStep(); current != nil && current.Resolve != nil {
		return current.Resolve(&m.Context, value)
	}

	return nil
}

// Helpers
