package portal

import (
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/cmdrun/domain/gameplay/dice"
	"solopg/app/cmdrun/domain/gameplay/session"
	"solopg/app/cmdrun/types/id"
	"solopg/app/shared/services/logs"
)

// - Portal - //
func Teleport() (*locations.Location, error) {
	list := locations.List()

	if len(list) == 0 {

		return nil, logs.Error("error.required", map[string]any{
			"Subject":  "location",
			"Property": "list",
		})
	}

	roll := dice.Roll(len(list))
	selectedLocation := list[roll-1]

	return &selectedLocation, nil
}

// Travel moves the session to the codex location entryID.
func Travel(gameSession *session.GameSession, entryID id.ID) {
	gameSession.UpdateLocation(gameSession.State.Codex.Locations.Find(entryID))
}
