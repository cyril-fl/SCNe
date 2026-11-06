package menumetadata

import (
	"fmt"
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters/races"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/shared/components/navigation"
	"scne/app/shared/services/i19n"
	"strings"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// NewSection reads the session on every render, so player and location changes show up live.
func NewSection(gameSession *session.GameSession) *navigation.Section {
	return navigation.NewSection(navigation.TemplateSection{
		ID: "metadata",
		Content: func() string {
			var content strings.Builder
			GetCharacterInfo(&content, gameSession)
			GetLocationInfo(&content, gameSession)
			GetWalletInfo(&content, gameSession)
			GetStatInfo(&content, gameSession)

			return strings.TrimRight(content.String(), "\n")
		},
	})
}

// GetCharacterInfo writes the player name, class and race lines to content.
func GetCharacterInfo(content *strings.Builder, gameSession *session.GameSession) {
	player := getState(gameSession).Player
	if player == nil {
		return
	}

	writeLine(content, "name", i19n.Localize(player.GetName()))
	writeLine(content, "class", plumtransform.Capitalize(i19n.Localize(player.Class)))
	writeLine(content, "race", plumtransform.Capitalize(i19n.Localize(races.NameOf(player.Race))))
}

// GetLocationInfo writes the current location line to content.
func GetLocationInfo(content *strings.Builder, gameSession *session.GameSession) {
	location := getState(gameSession).CurrentLocation
	if location == nil {
		writeLine(content, "location", plumtransform.Capitalize(i19n.Localize("unknown_place")))
		return
	}

	writeLine(content, "location", i19n.Localize(location.GetName()))
}

// GetWalletInfo writes the player wallet line to content.
func GetWalletInfo(content *strings.Builder, gameSession *session.GameSession) {
	player := getState(gameSession).Player
	if player == nil {
		return
	}

	writeLine(content, "wallet", i19n.Localize("wallet.entry", map[string]any{
		"Gold":   player.Wallet.Gold,
		"Silver": player.Wallet.Silver,
		"Copper": player.Wallet.Copper,
	}))
}

// GetStatInfo writes the player stats to content, two per line.
func GetStatInfo(content *strings.Builder, gameSession *session.GameSession) {
	content.WriteString("\n")
	content.WriteString(plumtransform.Capitalize(i19n.Localize("stats")))
	content.WriteString("\n")

	player := getState(gameSession).Player
	if player == nil {
		content.WriteString(i19n.Localize("stats.entry:unknown"))
		return
	}

	for i, stat := range stats.List() {
		if i%2 == 1 {
			content.WriteString("   ")
		} else if i > 0 {
			content.WriteString("\n")
		}
		label := plumtransform.Capitalize(i19n.Localize("stat." + string(stat)))
		fmt.Fprintf(content, "%-10s%5s", label, player.FormatStat(stat))
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

func writeLine(content *strings.Builder, labelID, value string) {
	fmt.Fprintf(content, "%-9s%s\n", plumtransform.Capitalize(i19n.Localize(labelID)), value)
}
