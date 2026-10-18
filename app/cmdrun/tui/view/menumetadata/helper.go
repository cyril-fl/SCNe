package menumetadata

import (
	"fmt"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/gameplay/session"
	"solopg/app/shared/services/i19n"
	"strings"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// REFACTOR LOW Surtout la maniere de if else
func GetCharacterInfo(content *strings.Builder, gameSession *session.GameSession) {
	var characterName = plumtransform.Uppercase(i19n.Localize("character"))

	isSessionNil := gameSession == nil || gameSession.State == nil

	if isSessionNil || gameSession.State.Player == nil {
		fmt.Fprintf(content, "%s", characterName)
	} else {
		fmt.Fprintf(content, "%s %s", characterName, i19n.Localize(gameSession.State.Player.GetName()))
	}

	content.WriteString("\n")
}

func GetLocationInfo(content *strings.Builder, gameSession *session.GameSession) {
	var locationName = plumtransform.Uppercase(i19n.Localize("location"))

	isSessionNil := gameSession == nil || gameSession.State == nil

	if isSessionNil || gameSession.State.CurrentLocation == nil {
		fmt.Fprintf(content, "%s %s", locationName, i19n.Localize("unknown_place"))
	} else {
		fmt.Fprintf(content, "%s %s", locationName, i19n.Localize(gameSession.State.CurrentLocation.GetName()))
	}
	content.WriteString("\n")
}

func GetStatInfo(content *strings.Builder, gameSession *session.GameSession) {
	content.WriteString(plumtransform.Uppercase(i19n.Localize("stats")))
	content.WriteString("\n")

	if gameSession == nil || gameSession.State == nil || gameSession.State.Player == nil {
		content.WriteString(i19n.Localize("stats.entry:unknowns"))
	} else {
		for _, stat := range stats.List() {
			key := "stat." + string(stat)
			label := i19n.Localize(key)
			fmt.Fprintf(content, "%-10s %d\n", label, gameSession.State.Player.Stats[stat])
		}
	}
}
