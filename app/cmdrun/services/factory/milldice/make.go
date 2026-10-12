package milldice

import (
	"solopg/app/cmdrun/domain/gameplay/dice"
	"solopg/app/cmdrun/services/factory/roll"
	"solopg/app/shared/services/factory"
	"solopg/app/shared/services/i19n"
)

type mill struct {
	factory.Mill
	dice dice.Dice
}

func New(value dice.Dice) roll.Roller {
	return &mill{dice: value}
}

func (m *mill) Roll() (string, error) {
	return i19n.Localize("dice.roll:result", map[string]any{
		"Dice":  m.dice.GetName(),
		"Value": m.dice.Roll(),
	}), nil
}
