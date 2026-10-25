package statemachine

import (
	plumprimitive "github.com/cyril-fl/plum/primitive"
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/layout/footer"
	"scne/app/shared/services/logs"
	interfass "scne/app/shared/types/interface"
	"slices"

	tea "charm.land/bubbletea/v2"
)

// TemplateViewState configures a ViewState: start state, allowed transitions and views.
type TemplateViewState[T comparable] struct {
	Base       T
	Transition map[T][]T
	Views      map[T]interfass.Viewable
}

// ViewState shows the view of its current state and only moves along allowed transitions.
type ViewState[T comparable] struct {
	plumprimitive.Fallible

	current    T
	transition map[T][]T
	view       map[T]interfass.Viewable
	size       size.Size
}

// NewView returns a ViewState in template.Base.
func NewView[T comparable](template TemplateViewState[T]) ViewState[T] {
	return ViewState[T]{
		current:    template.Base,
		transition: template.Transition,
		view:       template.Views,
	}
}

// Getters & Setters

// Current returns the current state.
func (s ViewState[T]) Current() T {
	return s.current
}

// GetView renders the current view, or "" without one.
func (s ViewState[T]) GetView() string {
	if view := s.view[s.current]; view == nil {
		return ""
	}

	return s.view[s.current].GetView()
}

// SetSize stores value and resizes the current view.
func (s *ViewState[T]) SetSize(value size.Size) {
	s.size = value
	if view := s.view[s.current]; view != nil {
		view.SetSize(value)
	}
}

// Methods

// Update forwards msg to the current view.
func (s *ViewState[T]) Update(msg tea.Msg) tea.Cmd {
	if view := s.view[s.current]; view == nil {
		return nil
	}

	return s.view[s.current].Update(msg)
}

// HandleTransition moves to state to and resizes its view.
// A disallowed transition stores an error and keeps the current state.
func (s *ViewState[T]) HandleTransition(to T) {
	if to == s.current {
		return
	}

	if !slices.Contains(s.transition[s.current], to) {
		s.SetErr(logs.MildError("error.invalid:transition"))
		return
	}

	s.current = to

	if view := s.view[to]; view != nil {
		view.SetSize(s.size)
	}
}

// Helpers

// Interface implementation

// Amend forwards content to the current view when it is Amandable.
func (s ViewState[T]) Amend(content string) {
	if view, ok := s.view[s.current].(interfass.Amandable); ok {
		view.Amend(content)
	}
}

// GetFooter returns the current view shortcuts.
func (s ViewState[T]) GetFooter() []string {
	var f []string
	if view, ok := s.view[s.current].(footer.HasFooter); ok {
		f = append(f, view.GetFooter()...)
	}
	return f
}
