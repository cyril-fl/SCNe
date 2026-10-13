package millroll

import (
	"solopg/app/cmdrun/domain/gameplay/dice"
	"solopg/app/cmdrun/domain/gameplay/hint"
	"solopg/app/cmdrun/domain/gameplay/oracle"

	"solopg/app/shared/services/factory"
	"solopg/app/shared/services/i19n"
	sharedtui "solopg/app/shared/tui"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Roller is an interface that defines a method for rolling game element and returning the result.
type Roller interface {
	Roll() (string, error)
}

// ResultMsg is a message type that contains the result of a roll operation.
type ResultMsg struct {
	Content string
}

// Oracles returns a Roller that rolls the given oracle and returns a localized message with the result.
func Oracles(value oracle.Oracle) Roller {
	return &milloracle{oracle: value}
}

type milloracle struct {
	factory.Mill
	oracle oracle.Oracle
}

// Hint returns a Roller that rolls the given hint and returns a localized message with the result.
func Hint(value hint.Hint) Roller {
	return &millhint{hint: value}
}

type millhint struct {
	factory.Mill
	hint hint.Hint
}

// Dice returns a Roller that rolls the given dice and returns a localized message with the result.
func Dice(value dice.Dice) Roller {
	return &milldice{dice: value}
}

type milldice struct {
	factory.Mill
	dice dice.Dice
}

// Roll returns a localized message for the given oracle, including the rolled value and critical status.
func (m *milloracle) Roll() (string, error) {
	result, err := oracle.Roll[any](m.oracle)
	if err != nil {
		return "", err
	}

	return i19n.Localize("oracles.roll:result", map[string]any{
		"Roll":     result.Roll,
		"Value":    result.Result,
		"Critical": result.Critical,
	}), nil
}

// Roll returns a localized message for the given hint, including all possible values.
func (m *millhint) Roll() (string, error) {
	values := hint.Roll(m.hint)
	localized := make([]string, 0, len(values))
	for _, value := range values {
		localized = append(localized, i19n.Localize(value))
	}

	return i19n.Localize("hint.roll:result", map[string]any{
		"Value": strings.Join(localized, ", "),
	}), nil
}

// Roll returns a localized message for the given dice, including the rolled value.
func (m *milldice) Roll() (string, error) {
	return i19n.Localize("dice.roll:result", map[string]any{
		"Dice":  m.dice.GetName(),
		"Value": m.dice.Roll(),
	}), nil
}

// Roll returns a localized message for the given oracle, including the rolled value and critical status.
func Execute(roller Roller) tea.Cmd {
	return func() tea.Msg {
		content, err := roller.Roll()
		if err != nil {
			return sharedtui.ErrorMsg{Err: err}
		}
		return ResultMsg{Content: content}
	}
}