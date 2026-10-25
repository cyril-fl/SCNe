package millstats

import (
	"errors"
	"os"
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters/classes"
	"scne/app/cmdrun/domain/card/characters/races"
	"scne/app/cmdrun/domain/gameplay/oracle"
	"scne/app/shared/services/logs"
	"strconv"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type mill struct {
	plumprimitive.Mill

	values    map[string]string
	modifiers []stats.Modifier
	oraclesId string
}

// New returns a stats generator without form values.
func New() *mill {
	return NewWithValues(nil)
}

// NewWithValues returns a stats generator reading values (stats, class, race, encounter).
func NewWithValues(values map[string]string) *mill {
	return &mill{
		values:    values,
		oraclesId: "stat_generation",
	}
}

// Getters & Setters

// GetStats returns the base stats with the generated modifiers applied.
func (m *mill) GetStats() stats.Stats {
	s := stats.GetBasic()
	s.ApplyModifiers(m.modifiers)
	return s
}

// GetModifiers returns the generated modifiers.
func (m *mill) GetModifiers() []stats.Modifier {
	return m.modifiers
}

// SetOracleId sets the oracle used for random modifiers.
func (m *mill) SetOracleId(entropy string) *mill {
	m.oraclesId = entropy
	return m
}

// Handlers
// GenerateFromValues generates full build for character, based on provided values (class, race, stats and fallback oracle).
func (m *mill) GenerateFromValues(params FromValuesParams) {
	m.reset()
	m.makeModifiersFromValuesWithFallback()
	m.makeModifiersFromAttributesValues()

	if params.Randomness {
		m.makeRandomModifiersFromOracleValue()
	}
}

// GenerateFromOracle replaces the modifiers with random ones from the oracle.
func (m *mill) GenerateFromOracle() {
	m.reset()
	m.makeRandomModifiersFromOracleValue()
}

func (m *mill) makeModifiersFromValuesWithFallback() {
	rules, err := oracle.GetByName(m.values["encounter"])
	if err != nil {

		m.SetErr(logs.Error("error.not_found:id-error", map[string]any{
			"Subject": "encounter",
			"ID":      m.values["encounter"],
			"Error":   err,
		}))
		return
	}

	for _, stat := range stats.List() {
		m.modifiers = append(m.modifiers, makeModifierWithRandomFallback(withFallbackTemplate{
			stat:     stat,
			value:    m.values[stat.String()],
			fallback: *rules,
		}))
	}
}

func (m *mill) makeModifiersFromAttributesValues() {
	m.modifiers = append(m.modifiers, GetBonus(m.values["class"], m.values["race"])...)
}

func (m *mill) makeRandomModifiersFromOracleValue() {
	modifiers, err := makeModifiersFromOracle(m.oraclesId)
	if err != nil {

		m.SetErr(logs.Error("error.invalid:new", map[string]any{
			"Subject": "stat",
			"Error":   err,
		}))
		return
	}

	m.modifiers = append(m.modifiers, modifiers...)
}

func (m *mill) reset() {
	m.modifiers = nil
	m.ClearErr()
}

// Helpers
// GetBonus returns the class and race (by ID) bonus, ignoring unknown ones.
func GetBonus(class, race string) []stats.Modifier {
	var bonus []stats.Modifier
	if c := classes.FindByName(class); c != nil {
		bonus = append(bonus, c.GetBonus()...)
	}
	if r := races.FindByID(race); r != nil {
		bonus = append(bonus, r.GetBonus()...)
	}
	return bonus
}

// FromValuesParams configures GenerateFromValues; Randomness adds oracle modifiers.
type FromValuesParams struct {
	Randomness bool
}

type withFallbackTemplate struct {
	stat     stats.Stat
	value    string
	fallback oracle.Oracle
}

func makeModifierWithRandomFallback(params withFallbackTemplate) stats.Modifier {
	if value, err := strconv.Atoi(params.value); err == nil && params.value != "" {
		return stats.Modifier{
			Stat:  params.stat,
			Value: value,
		}
	}

	roll, err := oracle.Roll[int](params.fallback)
	if err != nil {

		cwd, cwderr := os.Getwd()
		logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, err),
		})
		return stats.Modifier{
			Stat:  params.stat,
			Value: 0,
		}
	}

	return stats.Modifier{
		Stat:  params.stat,
		Value: roll.Result,
	}
}

func makeModifiersFromOracle(id string) ([]stats.Modifier, error) {
	var modifiers []stats.Modifier
	rules, err := oracle.GetByName(id)
	if err != nil {

		return nil, logs.Error("error.invalid:new: %w", map[string]any{
			"Subject": "stat",
			"Error":   err,
		})
	}

	for _, stat := range stats.List() {
		roll, err := oracle.Roll[int](*rules)
		if err != nil {
			cwd, cwderr := os.Getwd()
			return nil,
				logs.Error("error.unexpected", map[string]any{
					"Path":  cwd,
					"Error": errors.Join(cwderr, err),
				})
		}

		modifiers = append(modifiers, stats.Modifier{
			Stat:  stat,
			Value: roll.Result,
		})
	}

	return modifiers, nil
}
