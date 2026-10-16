package stats

import (
	"fmt"
	"slices"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"solopg/app/shared/types/primitive"
	"solopg/config"
	"strings"
)

// - Configuration & caching - //
type yamlConfig struct {
	Names     []Stat `yaml:"names"`
	BaseStats Stats  `yaml:"baseStats"`
}

var fileConfigPath = config.Current.Documents.Files.Stats

var cachedConfig *yamlConfig

func loadFromFile() error {
	params, err := yaml.LoadFromFile[yamlConfig](fileConfigPath)
	if err != nil {
		return err
	}

	cachedConfig = params

	return nil
}

// - Stats & Effects - //
type Stat string

type Stats map[Stat]int

type Modifier struct {
	Stat  Stat `yaml:"stat"`
	Value int  `yaml:"value"`
}
/* 
LOW : Pour quoi EFFECT a be soiin d'annotation, pas CHARECTR dans card.go
pourquoi :
type Effect struct {
	primitive.Descriptible `yaml:"description"`
	Modifier               Modifier `yaml:"modifier"`
}

et pas 

type Effect struct {
	primitive.Descriptible `yaml:"description"`
	Modifier `yaml:"modifier"`
}
*/
type Effect struct {
	primitive.Descriptible `yaml:"description"`
	Modifier               Modifier `yaml:"modifier"`
}

func New(mods []Modifier) Stats {
	s := make(Stats)
	s.ApplyModifiers(mods)
	return s
}

func GetBasic() Stats {
	if cachedConfig == nil {
		if err := loadFromFile(); err != nil {
			logs.SilentWarning("error.loading", map[string]any{
				"Subject": "stats",
				"Error":   err,
			})

			return nil
		}
	}

	copy := make(Stats)

	for stat, value := range cachedConfig.BaseStats {
		copy[stat] = value
	}

	return copy
}

func (s *Stats) ApplyModifier(mod Modifier) {
	if s == nil {
		return
	}

	if !mod.Stat.Validate() {
		valid := make([]string, 0, len(List()))
		for _, stat := range List() {
			valid = append(valid, i19n.Localize(stat.String()))
		}
		logs.SilentWarning("error.unexpected:value", map[string]any{
			"Subject":  "stat",
			"Expected": strings.Join(valid, ", "),
			"Value":    mod.Stat,
		})
		return
	}

	(*s)[mod.Stat] += mod.Value
}

func (s *Stats) ApplyModifiers(mods []Modifier) {
	for _, mod := range mods {
		s.ApplyModifier(mod)
	}
}

func (s Stat) Validate() bool {
	return slices.Contains(List(), s)
}

func (s Stats) String() string {
	if s == nil {
		return ""
	}

	result := make([]string, 0, len(s))
	for stat, value := range s {
		result = append(result, fmt.Sprintf("%s: %d", stat, value))
	}

	return fmt.Sprintf("{%s}", strings.Join(result, ", "))
}

func (s Stats) MappedString() map[string]any {
	m := make(map[string]any)
	for stat, value := range s {
		m[string(stat)] = value
	}
	return m
}

func (s Stat) String() string {
	return string(s)
}

func (m Modifier) String() string {
	return fmt.Sprintf("%s: %d", m.Stat, m.Value)
}

// - Collection - //
func List() []Stat {
	if cachedConfig == nil {
		if err := loadFromFile(); err != nil {
			logs.SilentWarning("error.loading:file", map[string]any{
				"Subject": "stats",
				"File":    fileConfigPath,
				"Error":   err,
			})
			return nil
		}
	}

	return cachedConfig.Names
}
