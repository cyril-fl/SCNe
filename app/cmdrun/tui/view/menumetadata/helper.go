package menumetadata

import (
	"fmt"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/gameplay/session"
	"solopg/app/shared/services/i19n"
	"strings"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

func GetCharacterInfo(content *strings.Builder, gameSession *session.GameSession) {
	var characterName = plumtransform.Uppercase(i19n.Localize("character"))

	state := getState(gameSession)
	isPlayerNil := state.Player == nil

	switch {
	case isPlayerNil:
		fmt.Fprintf(content, "%s", characterName)
	default:
		fmt.Fprintf(content, "%s %s", characterName, i19n.Localize(state.Player.GetName()))
	}

	content.WriteString("\n")
}

func GetLocationInfo(content *strings.Builder, gameSession *session.GameSession) {
	var locationName = plumtransform.Uppercase(i19n.Localize("location"))

	state := getState(gameSession)
	isLocationNil := state.CurrentLocation == nil

	switch {
	case isLocationNil:
		fmt.Fprintf(content, "%s %s", locationName, i19n.Localize("unknown_place"))
	default:
		fmt.Fprintf(content, "%s %s", locationName, i19n.Localize(state.CurrentLocation.GetName()))
	}

	content.WriteString("\n")
}

func GetStatInfo(content *strings.Builder, gameSession *session.GameSession) {
	content.WriteString(plumtransform.Uppercase(i19n.Localize("stats")))
	content.WriteString("\n")

	state := getState(gameSession)
	isPlayerNil := state.Player == nil

	switch {
	case isPlayerNil:
		content.WriteString(i19n.Localize("stats.entry:unknowns"))
	default:
		for _, stat := range stats.List() {
			key := "stat." + string(stat)
			label := i19n.Localize(key)
			fmt.Fprintf(content, "%-10s %d\n", label, state.Player.Stats[stat])
		}
	}
}

// Helpers

// getState returns an empty state when the session or its state is missing,
// so callers can read its fields without nil checks.
func getState(gameSession *session.GameSession) *session.SessionState {
	if gameSession == nil || gameSession.State == nil {
		return &session.SessionState{}
	}
	return gameSession.State
}
