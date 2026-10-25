package menuaction

import (
	"testing"

	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/gameplay/codex"
	"scne/app/cmdrun/domain/gameplay/session"
)

func TestFindOpponentsResolvesEntryIDsOnce(t *testing.T) {
	first := &characters.Character{}
	second := &characters.Character{}
	book := codex.New()
	book.Npcs.Add(first)
	book.Npcs.Add(second)
	gameSession := session.NewGameSession(book.Npcs.Entries[0].ID, &session.SessionState{Codex: book})

	firstID := book.Npcs.Entries[0].ID.String()
	values := map[string]string{"opponent.0": firstID, "opponent.1": firstID}

	got := findOpponents(gameSession, values)

	if len(got) != 1 || got[0] != first {
		t.Fatalf("findOpponents() = %v, want only the first NPC", got)
	}
}
