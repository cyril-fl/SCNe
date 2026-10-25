package stats

import (
	"fmt"
	"slices"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/config"
	"strings"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Names     []Stat `yaml:"names"`
	BaseStats Stats  `yaml:"baseStats"`
}

var fileConfigPath = config.GetConfig().Documents.Files.Stats

var cachedConfig = plumprimitive.NewCache[yamlConfig]()

func loadFromFile() (yamlConfig, error) {
	loader := plummillload.New(fileConfigPath)
	params := loader.LoadFromFile[yamlConfig]()
	if loader.HasErr() {
		return yamlConfig{}, logs.RawError("%w", loader.GetErr())
	}

	return *params, nil
}

// - Stats & Effects - //

// Stat is a stat name from the stats config file.
type Stat string

// Stats maps each stat to its value.
type Stats map[Stat]int

// Modifier adds Value to Stat.
type Modifier struct {
	Stat  Stat `yaml:"stat"`
	Value int  `yaml:"value"`
}

// Effect is a described Modifier.
type Effect struct {
	plumprimitive.Descriptible `yaml:"description"`
	Modifier                   Modifier `yaml:"modifier"`
}

// New returns stats built from mods.
func New(mods []Modifier) Stats {
	s := make(Stats)
	s.ApplyModifiers(mods)
	return s
}

// GetBasicWithError returns a copy of the configured base stats.
func GetBasicWithError() (Stats, error) {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}

	copy := make(Stats)

	for stat, value := range config.BaseStats {
		copy[stat] = value
	}

	return copy, nil
}

// GetBasic is GetBasicWithError that logs the load error once and returns nil on failure.
func GetBasic() Stats {
	values, err := GetBasicWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.SilentWarning("error.loading", map[string]any{
				"Subject": "stats",
				"Error":   err,
			})
		})
		return nil
	}

	return values
}

// ApplyModifier adds mod to s; an unknown stat is logged and ignored.
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

// ApplyModifiers applies each of mods.
func (s *Stats) ApplyModifiers(mods []Modifier) {
	for _, mod := range mods {
		s.ApplyModifier(mod)
	}
}

// RemoveModifiers reverts each of mods.
func (s *Stats) RemoveModifiers(mods []Modifier) {
	for _, mod := range mods {
		mod.Value = -mod.Value
		s.ApplyModifier(mod)
	}
}

// Validate reports whether s is a configured stat.
func (s Stat) Validate() bool {
	return slices.Contains(List(), s)
}

// String formats s as "{stat: value, ...}" in no particular order.
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

// MappedString returns s keyed by stat name, for templates.
func (s Stats) MappedString() map[string]any {
	m := make(map[string]any)
	for stat, value := range s {
		m[string(stat)] = value
	}
	return m
}

// String returns s as a string.
func (s Stat) String() string {
	return string(s)
}

// String formats m as "stat: value".
func (m Modifier) String() string {
	return fmt.Sprintf("%s: %d", m.Stat, m.Value)
}

// - Collection - //

// ListWithError returns a copy of the configured stat names; the file is loaded once.
func ListWithError() ([]Stat, error) {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}

	return append([]Stat(nil), config.Names...), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
func List() []Stat {
	values, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.SilentWarning("error.loading:file", map[string]any{
				"Subject": "stats",
				"File":    fileConfigPath,
				"Error":   err,
			})
		})
		return nil
	}

	return values
}
