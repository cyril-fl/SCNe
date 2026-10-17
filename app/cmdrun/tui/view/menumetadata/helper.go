package menumetadata

import (
	"fmt"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/services/game"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/utils/transform"
	"strings"
)

// REFACTOR LOW Surtout la maniere de if else
func GetCharacterInfo(content *strings.Builder, engine *game.Engine) {
	var characterName = transform.Uppercase(i19n.Localize("character"))

	isEngineNil := engine == nil || engine.State == nil

	if isEngineNil || engine.State.Player == nil {
		fmt.Fprintf(content, "%s", characterName)
	} else {
		fmt.Fprintf(content, "%s %s", characterName, i19n.Localize(engine.State.Player.GetName()))
	}

	content.WriteString("\n")
}

func GetLocationInfo(content *strings.Builder, engine *game.Engine) {
	var locationName = transform.Uppercase(i19n.Localize("location"))

	isEngineNil := engine == nil || engine.State == nil

	if isEngineNil || engine.State.CurrentLocation == nil {
		fmt.Fprintf(content, "%s %s", locationName, i19n.Localize("unknown_place"))
	} else {
		fmt.Fprintf(content, "%s %s", locationName, i19n.Localize(engine.State.CurrentLocation.GetName()))
	}
	content.WriteString("\n")
}

func GetStatInfo(content *strings.Builder, engine *game.Engine) {
	content.WriteString(transform.Uppercase(i19n.Localize("stats")))
	content.WriteString("\n")

	if engine == nil || engine.State == nil || engine.State.Player == nil {
		content.WriteString(i19n.Localize("stats.entry:unknowns"))
	} else {
		for _, stat := range stats.List() {
			key := "stat." + string(stat)
			label := i19n.Localize(key)
			fmt.Fprintf(content, "%-10s %d\n", label, engine.State.Player.Stats[stat])
		}
	}
}
