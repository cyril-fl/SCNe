package combatboard

import (
	"strings"
	"testing"

	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/chat"
	"scne/app/shared/components/form"

	tea "charm.land/bubbletea/v2"
)

func newTestModel(merged *[]string) (*model, *session.GameSession) {
	gameSession := session.NewGameSession([16]byte{}, &session.SessionState{
		Player:    &characters.Character{Stats: stats.Stats{"physical": 10}},
		Opponents: []*characters.Character{{Stats: stats.Stats{"physical": 8}}},
	})
	mainchat := chat.New(chat.Template{
		Register: func(_, message string) { *merged = append(*merged, message) },
	})

	m := NewModel(UiParams{Session: gameSession, Main: mainchat})
	m.SetSize(size.New(size.Template{Width: 40, Height: 20}))

	return m, gameSession
}

func TestSetSizeListsSessionOpponents(t *testing.T) {
	var merged []string
	m, _ := newTestModel(&merged)

	if got := len(m.opponents.GetList().Items()); got != 1 {
		t.Fatalf("opponents listed = %d, want 1", got)
	}
	if m.getTarget() == nil {
		t.Fatal("expected the first opponent as target")
	}

	if info := m.getTargetInfo(); !strings.Contains(info, " 8") {
		t.Fatalf("target info = %q, want the target stats", info)
	}
}

func TestStatShortcutRollsAgainstTarget(t *testing.T) {
	var merged []string
	m, _ := newTestModel(&merged)

	if cmd := m.Update(tea.KeyPressMsg{Code: '2', Mod: tea.ModAlt}); cmd == nil {
		t.Fatal("expected alt+2 to roll")
	}
	if _, ok := getStatShortcut("alt+99"); ok {
		t.Fatal("expected no stat bound to alt+99")
	}
}

func TestEscapeMergesThreadAndEndsCombat(t *testing.T) {
	var merged []string
	m, gameSession := newTestModel(&merged)
	m.Amend("I strike")

	cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})

	if len(merged) != 2 || merged[1] != "I strike" {
		t.Fatalf("merged = %v, want the header then the combat line", merged)
	}
	if gameSession.State.Opponents != nil {
		t.Fatal("expected the combat to be ended")
	}
	if _, ok := cmd().(form.Close); !ok {
		t.Fatal("expected form.Close")
	}
	if len(m.buffer) != 0 {
		t.Fatal("expected an empty buffer after merge")
	}
}

func TestStatShortcutAcceptsConfiguredAlias(t *testing.T) {
	alt, okAlt := getStatShortcut("alt+2")
	alias, okAlias := getStatShortcut("„")

	if !okAlt || !okAlias || alt != alias {
		t.Fatalf("alt+2 = %q (%v), „ = %q (%v), want the same stat", alt, okAlt, alias, okAlias)
	}
}

func TestStatShortcutNeedsTarget(t *testing.T) {
	var merged []string
	m, gameSession := newTestModel(&merged)
	gameSession.State.Opponents = nil
	m.SetSize(size.New(size.Template{Width: 40, Height: 20}))

	if cmd := m.Update(tea.KeyPressMsg{Code: '2', Mod: tea.ModAlt}); cmd != nil {
		t.Fatal("expected no roll without a target")
	}
}

func TestOracleSectionFooterShowsStatShortcuts(t *testing.T) {
	var merged []string
	m, _ := newTestModel(&merged)

	footer := strings.Join(m.GetOracleSection().GetFooter(), " ")
	if !strings.Contains(footer, "⌥2") {
		t.Fatalf("oracle footer = %q, want the stat shortcuts", footer)
	}
}
