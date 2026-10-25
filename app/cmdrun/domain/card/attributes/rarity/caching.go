package rarity

import (
	"errors"
	"io/fs"
	"scne/app/shared/services/logs"
	"scne/config"
	"slices"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Values  []Rarity `yaml:"values"`
	Default Rarity   `yaml:"default"`
}

var fileConfigPath = config.GetConfig().Documents.Files.Rarity

var cachedConfig = plumprimitive.NewCache[yamlConfig]()

func loadFromFile() (yamlConfig, error) {
	// rarity is optional: no configured or existing file means no values, not an error.
	if fileConfigPath == "" {
		return yamlConfig{}, nil
	}

	loader := plummillload.New(fileConfigPath)
	params := loader.LoadFromFile[yamlConfig]()
	if errors.Is(loader.GetErr(), fs.ErrNotExist) {
		return yamlConfig{}, nil
	}
	if loader.HasErr() {
		return yamlConfig{}, logs.RawError("%w", loader.GetErr())
	}

	return *params, nil
}

func (c *yamlConfig) addValue(value Rarity) {
	if !slices.Contains(c.Values, value) {
		c.Values = append(c.Values, value)
	}
}

// - Rarity - //

// Rarity is a value from the rarity config file.
type Rarity string

// ListWithError returns a copy of the configured values; the file is loaded once.
func ListWithError() ([]Rarity, error) {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}

	return append([]Rarity(nil), config.Values...), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
func List() []Rarity {
	values, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.Error("error.loading:file", map[string]any{
				"Subject": "rarity",
				"File":    fileConfigPath,
				"Error":   err,
			})
		})
		return nil
	}

	return values
}

// Default returns the configured default rarity, or registers "" when unset or invalid.
func Default() Rarity {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return MakeDefault("")
	}

	if config.Default == "" || !config.Default.Validate() {
		return MakeDefault("")
	}

	return config.Default
}

// Validate reports whether r is a configured value.
func (r Rarity) Validate() bool {
	return slices.Contains(List(), r)
}

// MakeDefault adds defaultValue to the cached values when missing and returns it.
func MakeDefault(defaultValue string) Rarity {
	_, _ = ListWithError()
	cachedConfig.Update(func(config *yamlConfig) {
		config.addValue(Rarity(defaultValue))
	})
	return Rarity(defaultValue)
}

// AssertWithDefault returns provided when valid, otherwise Default.
func AssertWithDefault(provided string) Rarity {
	if Assert(provided) {
		return Rarity(provided)
	}

	return Default()
}

// Assert reports whether provided is a configured value.
func Assert(provided string) bool {
	return Rarity(provided).Validate()
}

// String returns r as a string.
func (r Rarity) String() string {
	return string(r)
}
