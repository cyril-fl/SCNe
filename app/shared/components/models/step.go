package models

import (
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/types/step"

	tea "charm.land/bubbletea/v2"
)

type ContextStepList[T any] = step.List[T]
type Step[T any] = step.Step[T]

type StepModel[T any] struct {
	*ContextStepList[T]
	Context T
	Size *tea.WindowSizeMsg
}

func NewStepModel[T any](steps *ContextStepList[T]) StepModel[T] {
	return StepModel[T]{
		ContextStepList: steps,
	}
}

func (sm *StepModel[T]) ResolveCurrent(value any) error {
	if current := sm.CurrentStep(); current != nil && current.Resolve != nil {
		return current.Resolve(&sm.Context, value)
	}

	return nil
}

func (sm *StepModel[T]) SelectNext() bool {
	for {
		next := sm.NextStep()
		if next == nil {
			return false
		}

		if next.Skip == nil || !next.Skip(&sm.Context) {
			return true
		}
	}
}

func (sm *StepModel[T]) ForwardNext() (tea.Cmd) {
	if !sm.SelectNext() {
		return  tea.Quit
	}

	submodel := sm.GetCurrentSubmodel()
	if submodel == nil {
		return tea.Quit
	}

	initCmd := submodel.Init()
	if sm.Size == nil {
		return initCmd
	}

	return tea.Sequence(
		initCmd, 
		sharedtui.SendSizeMsg(sm.Size),
	)
}

// Helpers 
