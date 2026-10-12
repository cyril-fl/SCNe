package milloracle

import (
	"solopg/app/cmdrun/domain/gameplay/oracle"
	"solopg/app/cmdrun/services/factory/roll"
	"solopg/app/shared/services/factory"
	"solopg/app/shared/services/i19n"
)

type mill struct {
	factory.Mill
	oracle oracle.Oracle
}

func New(value oracle.Oracle) roll.Roller {
	return &mill{oracle: value}
}

func (m *mill) Roll() (string, error) {
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
