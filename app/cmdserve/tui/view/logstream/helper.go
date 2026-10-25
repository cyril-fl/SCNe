package logstream

import (
	"scne/app/shared/services/logs"
)

func (m *model[T]) assertFilterFlag() {
	if filter, err := m.flags.GetString("filter"); err != nil {
	} else if filter == "" {
	} else if isKind := logs.AssertKind(logs.Kind(filter)); isKind {
	} else {
		logs.Error("error.unexpected:action", map[string]any{
			"Action": "unexpected:action.logs:filter",
			"Error": logs.CeaseError("error.invalid", map[string]any{
				"Subject": "filter",
			}),
		})
	}
}
