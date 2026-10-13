package statemachine

import (
	"slices"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/layout/footer"
	"solopg/app/shared/services/logs"
	interfass "solopg/app/shared/types/interface"
	"solopg/app/shared/types/primitive"

	tea "charm.land/bubbletea/v2"
)

type TemplateViewState[T comparable] struct {
	Base       T
	Transition map[T][]T
	Views      map[T]interfass.Viewable
}

type ViewState[T comparable] struct {
	primitive.Fallible

	current    T
	transition map[T][]T
	view       map[T]interfass.Viewable
	size       size.Size
}

func NewView[T comparable](template TemplateViewState[T]) ViewState[T] {
	return ViewState[T]{
		current:    template.Base,
		transition: template.Transition,
		view:       template.Views,
	}
}

// Getters & Setters
func (s ViewState[T]) Current() T {
	return s.current
}

func (s ViewState[T]) GetView() string {
	if view := s.view[s.current]; view == nil {
		return ""
	}

	return s.view[s.current].GetView()
}

func (s *ViewState[T]) SetSize(value size.Size) {
	s.size = value
	if view := s.view[s.current]; view != nil {
		view.SetSize(value)
	}
}

// Methods
func (s *ViewState[T]) Update(msg tea.Msg) tea.Cmd {
	if view := s.view[s.current]; view == nil {
		return nil
	}

	return s.view[s.current].Update(msg)
}

func (s *ViewState[T]) HandleTransition(to T) {
	if to == s.current {
		return
	}

	if !slices.Contains(s.transition[s.current], to) {
		s.SetErr(logs.MildError("Transition pas autorise"))
		return
	}

	s.current = to

	if view := s.view[to]; view != nil {
		view.SetSize(s.size)
	}
}

// Helpers

// Interface implementation
func (s ViewState[T]) Amend(content string) {
	if view, ok := s.view[s.current].(interfass.Amandable); ok {
		view.Amend(content)
	}
}

func (s ViewState[T]) GetFooter() []string {
	var f []string
	if view, ok := s.view[s.current].(footer.HasFooter); ok {
		f = append(f, view.GetFooter()...)
	}
	return f
}
