package models

import (
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/utils/transform"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type RerollModel[T any] struct {
	Attempt int
	Limit   int
	Choice  list.Model
	Value   T
	Roll    func() (T, error)
}

func NewRerollModel[T any](options list.Model, handleRoll func() (T, error), limit int) RerollModel[T] {
	value, err := handleRoll()

	if err != nil {
		panic(logs.Error("error.unexpected:action", map[string]any{
			"Action": "unexpected:action.roll",
			"Error":  err,
		}))
	}

	return RerollModel[T]{
		Attempt: 1,
		Limit:   limit,
		Choice:  options,
		Value:   value,
		Roll:    handleRoll,
	}
}

func (s *RerollModel[T]) Reroll() {
	value, err := s.Roll()

	if err != nil {
		panic(logs.Error("error.unexpected:action", map[string]any{
			"Action": "unexpected:action.roll",
			"Error":  err,
		}))
	}

	s.Attempt++
	s.Value = value

	s.checkAttemptLimit()
}

func (s *RerollModel[T]) checkAttemptLimit() {
	if s.Attempt < s.Limit {
		return
	}

	s.Choice = makeConfirmationModel(s.Choice)
}

func (s *RerollModel[T]) IsOutOfLimit() bool {
	return s.Attempt >= s.Limit
}

func (s *RerollModel[T]) HandleEnterInput() tea.Cmd {
	isSelected, ok := s.Choice.SelectedItem().(Item[bool])
	if !ok {
		return nil
	}

	if isSelected.Value() || s.IsOutOfLimit() {
		return sharedtui.SendResolutionMsg(sharedtui.ResolutionMsg{Completed: true, Value: s.Value})
	}

	s.Reroll()

	return nil
}

func NewOptionsModel(options rerollOptions) list.Model {
	items := []list.Item{
		NewItem(trueLabel(options), "", true),
		NewItem(falseLabel(options), "", false),
	}

	return newConfiguredList(items, 0, 0)
}

func makeConfirmationModel(m list.Model) list.Model {
	items := []list.Item{
		NewItem(trueLabel(DefaultRerollOptions), "", true),
	}

	return newConfiguredList(items, m.Width(), m.Height())
}

func newConfiguredList(items []list.Item, width, height int) list.Model {
	model := list.New(items, list.NewDefaultDelegate(), width, height)
	ConfigureList(&model)
	return model
}

func trueLabel(options rerollOptions) string {
	return transform.Capitalize(i19n.Localize(options.true))
}

func falseLabel(options rerollOptions) string {
	return transform.Capitalize(i19n.Localize(options.false))
}

type rerollOptions struct {
	true  string
	false string
}

var DefaultRerollOptions = rerollOptions{
	true:  "accept",
	false: "reroll",
}
