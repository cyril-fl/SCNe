package millroll

import (
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/gameplay/dice"
	"scne/app/cmdrun/domain/gameplay/hint"
	"scne/app/cmdrun/domain/gameplay/oracle"

	"scne/app/shared/messages"

	plumprimitive "github.com/cyril-fl/plum/primitive"

	"scne/app/shared/services/i19n"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Roller is an interface that defines a method for rolling game element and returning the result.
type Roller interface {
	Roll() (string, error)
}

// Oracles returns a Roller that rolls the given oracle and returns a localized message with the result.
func Oracles(value oracle.Oracle) Roller {
	return &milloracle{oracle: value}
}

type milloracle struct {
	plumprimitive.Mill
	oracle oracle.Oracle
}

// Hint returns a Roller that rolls the given hint and returns a localized message with the result.
func Hint(value hint.Hint) Roller {
	return &millhint{hint: value}
}

type millhint struct {
	plumprimitive.Mill
	hint hint.Hint
}

// Dice returns a Roller that rolls the given dice and returns a localized message with the result.
func Dice(value dice.Dice) Roller {
	return &milldice{dice: value}
}

type milldice struct {
	plumprimitive.Mill
	dice dice.Dice
}

// Fight returns a Roller that rolls the given oracle against target on stat and returns a localized message with the result.
func Fight(target *characters.Character, stat stats.Stat, value oracle.Oracle) Roller {
	return &millfight{target: target, stat: stat, oracle: value}
}

type millfight struct {
	plumprimitive.Mill
	target *characters.Character
	stat   stats.Stat
	oracle oracle.Oracle
}

// Roll rolls the oracle and returns a localized message naming the stat, the target and the result.
func (m *millfight) Roll() (string, error) {
	result, err := oracle.Roll[any](m.oracle)
	if err != nil {
		return "", err
	}

	return i19n.Localize("combat.roll:result", map[string]any{
		"Stat":     "stat." + m.stat.String(),
		"Opponent": m.target.GetName(),
		"Oracle":   m.oracle.GetName(),
		"Roll":     result.Roll,
		"Value":    result.Result,
	}), nil
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

// Execute returns a command emitting the roll result, or messages.Failed on error.
func Execute(roller Roller) tea.Cmd {
	return func() tea.Msg {
		content, err := roller.Roll()
		if err != nil {
			return messages.Failed{Err: err}
		}
		return messages.Result{Content: content}
	}
}
