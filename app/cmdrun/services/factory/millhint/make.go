package millhint

import (
	"strings"

	"solopg/app/cmdrun/domain/gameplay/hint"
	"solopg/app/cmdrun/services/factory/roll"
	"solopg/app/shared/services/factory"
	"solopg/app/shared/services/i19n"
)

type mill struct {
	factory.Mill
	hint hint.Hint
}

func New(value hint.Hint) roll.Roller {
	return &mill{hint: value}
}

func (m *mill) Roll() (string, error) {
	values := hint.Roll(m.hint)
	localized := make([]string, 0, len(values))
	for _, value := range values {
		localized = append(localized, i19n.Localize(value))
	}

	return i19n.Localize("hint.roll:result", map[string]any{
		"Value": strings.Join(localized, ", "),
	}), nil
}
