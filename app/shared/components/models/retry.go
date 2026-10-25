package models

import (
	"solopg/app/shared/messages"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// RerollModel holds a rolled value the user can accept or reroll up to Limit attempts.
type RerollModel[T any] struct {
	Attempt int
	Limit   int
	Choice  list.Model
	Value   T
	Roll    func() (T, error)
}

// NewRerollModel rolls the first value and panics if handleRoll fails.
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

// Reroll rolls a new value and panics if the roll fails.
// Once the limit is reached, only the accept option remains.
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

// IsOutOfLimit reports whether no reroll is left.
func (s *RerollModel[T]) IsOutOfLimit() bool {
	return s.Attempt >= s.Limit
}

// HandleEnterInput resolves with the value when accepted or out of rerolls, otherwise rerolls.
func (s *RerollModel[T]) HandleEnterInput() tea.Cmd {
	isSelected, ok := s.Choice.SelectedItem().(Item[bool])
	if !ok {
		return nil
	}

	if isSelected.Value() || s.IsOutOfLimit() {
		return messages.Resolve(messages.Resolved{Completed: true, Value: s.Value})
	}

	s.Reroll()

	return nil
}

// NewOptionsModel returns the accept/reroll choice list.
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
	return plumtransform.Capitalize(i19n.Localize(options.true))
}

func falseLabel(options rerollOptions) string {
	return plumtransform.Capitalize(i19n.Localize(options.false))
}

type rerollOptions struct {
	true  string
	false string
}

// DefaultRerollOptions holds the accept and reroll translation keys.
var DefaultRerollOptions = rerollOptions{
	true:  "accept",
	false: "reroll",
}
