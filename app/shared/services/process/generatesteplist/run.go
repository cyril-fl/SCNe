package generatesteplist

import (
	"scne/app/shared/services/logs"
	"scne/app/shared/types/step"

	tea "charm.land/bubbletea/v2"
	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type run[T any] struct {
	plumprimitive.Process

	steps step.List[T]
}

// Generator builds a step list and runs it as a tea program.
type Generator[T any] = run[T]

// Process returns an empty Generator.
func Process[T any]() *run[T] {
	return &run[T]{}
}

// Run runs model as a tea program, storing any error.
func (p *run[T]) Run(model tea.Model, opts ...tea.ProgramOption) {
	program := tea.NewProgram(model, opts...)
	_, err := program.Run()

	p.SetErr(err)
}

// Getters & Setters

// GetResult is not implemented; it only logs a warning.
func (p *run[T]) GetResult() {
	logs.SilentWarning("error:not_implemented", map[string]any{
		"Function": "GetResult",
		"Subject":  "session",
	})
}

// GetSteps returns the step list.
func (p *run[T]) GetSteps() *step.List[T] {
	return &p.steps
}

// Methods

// Add appends steps.
func (p *run[T]) Add(steps ...step.Step[T]) *run[T] {
	p.steps.Steps = append(p.steps.Steps, steps...)
	return p
}

// Insert adds steps right after the current one.
func (p *run[T]) Insert(steps ...step.Step[T]) *run[T] {
	if len(steps) == 0 {
		return p
	}

	insertAt := p.steps.CurrentIndex + 1
	if insertAt < 0 {
		insertAt = 0
	}
	if insertAt > len(p.steps.Steps) {
		insertAt = len(p.steps.Steps)
	}

	oldSteps := p.steps.Steps
	newSteps := make([]step.Step[T], 0, len(oldSteps)+len(steps))

	newSteps = append(newSteps, oldSteps[:insertAt]...)
	newSteps = append(newSteps, steps...)
	newSteps = append(newSteps, oldSteps[insertAt:]...)

	p.steps.Steps = newSteps
	return p
}

// Helpers
