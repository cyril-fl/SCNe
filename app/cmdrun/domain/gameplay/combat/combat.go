package combat

import (
	"maps"

	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/characters"
	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/gameplay/oracle"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
)

// shortcut: health stat hardcoded (stats.yaml marks it as HP); move to config if stats become configurable per system.
const healthStat stats.Stat = "health"

// - Fighter - //

// Fighter is a character in a combat with the actions of its class.
type Fighter struct {
	Character *characters.Character
	Actions   []Action
	Intent    *oracle.Oracle
	guard     int
}

// GetHealth returns the fighter health.
func (f Fighter) GetHealth() int {
	return f.Character.Stats[healthStat]
}

func (f Fighter) findAction(id string) *Action {
	for _, action := range f.Actions {
		if action.id == id {
			return &action
		}
	}
	return nil
}

func newFighter(character *characters.Character) (*Fighter, error) {
	if character == nil {
		return nil, logs.MildError("error.required", map[string]any{
			"Subject":  "combat",
			"Property": "fighter",
		})
	}

	if character.Stats[healthStat] <= 0 {
		return nil, logs.MildError("combat.error:exhausted", map[string]any{
			"Fighter": character.GetName(),
		})
	}

	class := classes.FindByName(character.Class)
	if class == nil || len(class.GetActions()) == 0 {
		return nil, logs.MildError("combat.error:no_action", map[string]any{
			"Fighter": character.GetName(),
		})
	}

	actions := make([]Action, 0, len(class.GetActions()))
	for _, id := range class.GetActions() {
		action, err := GetByID(id)
		if err != nil {
			return nil, err
		}
		actions = append(actions, *action)
	}

	fighter := &Fighter{
		Character: character,
		Actions:   actions,
	}

	if class.GetIntent() == "" {
		return fighter, nil
	}

	intent, err := oracle.GetByName(class.GetIntent())
	if err != nil {
		return nil, err
	}
	fighter.Intent = intent

	return fighter, nil
}

// - Combat - //

// Combat is a running exchange between the player and an opponent.
// Damage is applied to the player character; the opponent is a copy.
type Combat struct {
	Player   *Fighter
	Opponent *Fighter
	Round    int
	roll     func(oracle.Oracle) (*oracle.Result[string], error)
}

type move struct {
	fighter *Fighter
	target  *Fighter
	effect  map[string]int
}

// New returns a combat between player and a copy of opponent.
// It fails when a fighter has no health, no class action, or the opponent has no intent.
func New(player, opponent *characters.Character) (*Combat, error) {
	playerFighter, err := newFighter(player)
	if err != nil {
		return nil, err
	}

	var copied *characters.Character
	if opponent != nil {
		clone := *opponent
		clone.Stats = maps.Clone(opponent.Stats)
		copied = &clone
	}

	opponentFighter, err := newFighter(copied)
	if err != nil {
		return nil, err
	}

	if opponentFighter.Intent == nil {
		return nil, logs.MildError("combat.error:no_intent", map[string]any{
			"Fighter": opponent.GetName(),
		})
	}

	return &Combat{
		Player:   playerFighter,
		Opponent: opponentFighter,
		roll:     oracle.Roll[string],
	}, nil
}

// Getters & Setters

// IsOver reports whether a fighter has no health left.
func (c Combat) IsOver() bool {
	return c.Player.GetHealth() <= 0 || c.Opponent.GetHealth() <= 0
}

// GetOutcome returns the end translation key, or "" while the combat runs.
func (c Combat) GetOutcome() string {
	isPlayerDown := c.Player.GetHealth() <= 0
	isOpponentDown := c.Opponent.GetHealth() <= 0

	switch {
	case isPlayerDown && isOpponentDown:
		return "combat.end:draw"
	case isPlayerDown:
		return "combat.end:defeat"
	case isOpponentDown:
		return "combat.end:victory"
	}
	return ""
}

// Methods

// Play runs one round: the player plays actionID, the opponent plays its rolled intent.
// Every roll happens before any effect, so a failed roll leaves the combat untouched.
// Guards apply before damage. It returns the localized round events.
func (c *Combat) Play(actionID string) ([]string, error) {
	if c.IsOver() {
		return nil, logs.MildError("combat.error:over")
	}

	playerAction := c.Player.findAction(actionID)
	if playerAction == nil {
		return nil, logs.MildError("combat.error:unavailable", map[string]any{
			"Action": "actions." + actionID,
		})
	}

	events := []string{i19n.Localize("combat.round", map[string]any{"Round": c.Round + 1})}

	playerMove, event, err := c.resolve(c.Player, c.Opponent, playerAction)
	if err != nil {
		return nil, err
	}
	events = append(events, event)

	opponentAction, err := c.rollIntent()
	if err != nil {
		return nil, err
	}

	moves := []move{playerMove}

	if opponentAction == nil {
		events = append(events, i19n.Localize("combat.hesitate", map[string]any{
			"Fighter": c.Opponent.Character.GetName(),
		}))
	} else {
		opponentMove, event, err := c.resolve(c.Opponent, c.Player, opponentAction)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
		moves = append(moves, opponentMove)
	}

	for _, m := range moves {
		m.fighter.guard += m.effect[effectGuard]
	}

	for _, m := range moves {
		events = append(events, c.applyDamage(m)...)
	}

	c.Player.guard = 0
	c.Opponent.guard = 0
	c.Round++

	if c.IsOver() {
		events = append(events, i19n.Localize(c.GetOutcome()))
	}

	return events, nil
}

// Helpers

func (c *Combat) resolve(fighter, target *Fighter, action *Action) (move, string, error) {
	actionOracle, err := oracle.GetByName(action.oracle)
	if err != nil {
		return move{}, "", err
	}

	result, err := c.roll(*actionOracle)
	if err != nil {
		return move{}, "", err
	}

	event := i19n.Localize("combat.action:result", map[string]any{
		"Fighter": fighter.Character.GetName(),
		"Action":  action.GetName(),
		"Roll":    result.Roll,
		"Outcome": "combat.outcome." + result.Result,
	})

	return move{
		fighter: fighter,
		target:  target,
		effect:  action.effects[result.Result],
	}, event, nil
}

// rollIntent returns the opponent action rolled on its intent oracle, or nil when it has no such action.
func (c *Combat) rollIntent() (*Action, error) {
	result, err := c.roll(*c.Opponent.Intent)
	if err != nil {
		return nil, err
	}
	return c.Opponent.findAction(result.Result), nil
}

func (c *Combat) applyDamage(m move) []string {
	if m.effect[effectDamage] == 0 {
		return nil
	}

	damage := max(0, m.effect[effectDamage]-m.target.guard)

	m.target.Character.Stats.ApplyModifier(stats.Modifier{
		Stat:  healthStat,
		Value: -damage,
	})

	return []string{i19n.Localize("combat.damage", map[string]any{
		"Fighter": m.target.Character.GetName(),
		"Damage":  damage,
		"Health":  m.target.GetHealth(),
	})}
}
