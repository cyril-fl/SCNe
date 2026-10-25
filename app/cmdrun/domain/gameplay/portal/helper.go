package portal

import (
	"scne/app/cmdrun/domain/card/locations"
	"scne/app/cmdrun/domain/gameplay/dice"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/logs"
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
