package combat

import (
	"slices"
	"testing"

	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/characters"
	"solopg/app/cmdrun/domain/gameplay/oracle"
	"solopg/app/shared/services/logs"
)

type combatTestLogRegistry struct{}

func (combatTestLogRegistry) Register(*logs.Log) error { return nil }

func newTestCharacter(class string, health int) *characters.Character {
	return &characters.Character{
		Class: class,
		Stats: stats.Stats{healthStat: health},
	}
}

// newTestCombat returns a warrior vs warrior combat whose rolls return results in order:
// player action, opponent intent, opponent action.
func newTestCombat(t *testing.T, playerHealth, opponentHealth int, results ...string) (*Combat, *characters.Character, *characters.Character) {
	t.Helper()
	t.Cleanup(func() { logs.Init(nil) })
	logs.Init(combatTestLogRegistry{})

	player := newTestCharacter("class.warrior", playerHealth)
	opponent := newTestCharacter("class.warrior", opponentHealth)

	c, err := New(player, opponent)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	c.roll = func(oracle.Oracle) (*oracle.Result[string], error) {
		if len(results) == 0 {
			t.Fatal("unexpected roll")
		}
		result := results[0]
		results = results[1:]
		return &oracle.Result[string]{Roll: 10, Result: result}, nil
	}

	return c, player, opponent
}

func TestPlayHitDamagesPlayerAndOpponentCopy(t *testing.T) {
	c, player, opponent := newTestCombat(t, 10, 10, "hit", "attack", "hit")

	if _, err := c.Play("attack"); err != nil {
		t.Fatalf("Play() error = %v", err)
	}

	if player.Stats[healthStat] != 7 {
		t.Fatalf("player health = %d, want 7", player.Stats[healthStat])
	}
	if c.Opponent.GetHealth() != 7 {
		t.Fatalf("opponent health = %d, want 7", c.Opponent.GetHealth())
	}
	if opponent.Stats[healthStat] != 10 {
		t.Fatalf("codex opponent health = %d, want 10 (untouched)", opponent.Stats[healthStat])
	}
	if c.Round != 1 {
		t.Fatalf("Round = %d, want 1", c.Round)
	}
}

func TestPlayGuardReducesCriticalDamage(t *testing.T) {
	c, player, _ := newTestCombat(t, 10, 10, "critical", "defend", "success")

	if _, err := c.Play("attack"); err != nil {
		t.Fatalf("Play() error = %v", err)
	}

	if c.Opponent.GetHealth() != 7 {
		t.Fatalf("opponent health = %d, want 7 (6 critical - 3 guard)", c.Opponent.GetHealth())
	}
	if player.Stats[healthStat] != 10 {
		t.Fatalf("player health = %d, want 10", player.Stats[healthStat])
	}
}

func TestPlayGuardNeverHeals(t *testing.T) {
	c, player, _ := newTestCombat(t, 10, 10, "critical", "attack", "hit")

	if _, err := c.Play("defend"); err != nil {
		t.Fatalf("Play() error = %v", err)
	}

	if player.Stats[healthStat] != 10 {
		t.Fatalf("player health = %d, want 10 (3 damage - 6 guard floors at 0)", player.Stats[healthStat])
	}
}

func TestPlayGuardLastsOneRound(t *testing.T) {
	c, player, _ := newTestCombat(t, 10, 10, "success", "attack", "miss", "miss", "attack", "hit")

	if _, err := c.Play("defend"); err != nil {
		t.Fatalf("first Play() error = %v", err)
	}
	if _, err := c.Play("attack"); err != nil {
		t.Fatalf("second Play() error = %v", err)
	}

	if player.Stats[healthStat] != 7 {
		t.Fatalf("player health = %d, want 7", player.Stats[healthStat])
	}
}

func TestPlayRejectsActionOutsideClass(t *testing.T) {
	c, _, _ := newTestCombat(t, 10, 10)

	if _, err := c.Play("fireball"); err == nil {
		t.Fatal("expected an error for an action outside the class")
	}
	if c.Round != 0 {
		t.Fatalf("Round = %d, want 0", c.Round)
	}
}

func TestPlayOpponentHesitatesOnUnknownIntent(t *testing.T) {
	c, player, _ := newTestCombat(t, 10, 10, "miss", "dance")

	events, err := c.Play("attack")
	if err != nil {
		t.Fatalf("Play() error = %v", err)
	}

	if player.Stats[healthStat] != 10 {
		t.Fatalf("player health = %d, want 10", player.Stats[healthStat])
	}
	if !slices.ContainsFunc(events, func(event string) bool { return event == "{{combat.hesitate}}" || event == "combat.hesitate" }) {
		t.Fatalf("events = %v, want a hesitate event", events)
	}
}

func TestPlayEndsCombat(t *testing.T) {
	tests := []struct {
		name           string
		playerHealth   int
		opponentHealth int
		results        []string
		want           string
	}{
		{name: "victory", playerHealth: 10, opponentHealth: 3, results: []string{"hit", "attack", "miss"}, want: "combat.end:victory"},
		{name: "defeat", playerHealth: 3, opponentHealth: 10, results: []string{"miss", "attack", "hit"}, want: "combat.end:defeat"},
		{name: "draw", playerHealth: 3, opponentHealth: 3, results: []string{"hit", "attack", "hit"}, want: "combat.end:draw"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, _, _ := newTestCombat(t, test.playerHealth, test.opponentHealth, test.results...)

			if _, err := c.Play("attack"); err != nil {
				t.Fatalf("Play() error = %v", err)
			}

			if !c.IsOver() || c.GetOutcome() != test.want {
				t.Fatalf("IsOver() = %v, GetOutcome() = %q, want %q", c.IsOver(), c.GetOutcome(), test.want)
			}
			if _, err := c.Play("attack"); err == nil {
				t.Fatal("expected an error when playing a finished combat")
			}
		})
	}
}

func TestNewRejectsUnfitFighters(t *testing.T) {
	t.Cleanup(func() { logs.Init(nil) })
	logs.Init(combatTestLogRegistry{})

	tests := []struct {
		name     string
		player   *characters.Character
		opponent *characters.Character
	}{
		{name: "player without health", player: newTestCharacter("class.warrior", 0), opponent: newTestCharacter("class.warrior", 10)},
		{name: "opponent without class", player: newTestCharacter("class.warrior", 10), opponent: newTestCharacter("", 10)},
		{name: "missing opponent", player: newTestCharacter("class.warrior", 10), opponent: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := New(test.player, test.opponent); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
