package variety

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
	Values  []Variety `yaml:"values"`
	Default Variety   `yaml:"default"`
}

var fileConfigPath = config.GetConfig().Documents.Files.Variety

var cachedConfig = plumprimitive.NewCache[yamlConfig]()

func loadFromFile() (yamlConfig, error) {
	// variety is optional: no configured or existing file means no values, not an error.
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

func (c *yamlConfig) addValue(value Variety) {
	if !slices.Contains(c.Values, value) {
		c.Values = append(c.Values, value)
	}
}

// - Variety - //

// Variety is a value from the variety config file.
type Variety string

// ListWithError returns a copy of the configured values; the file is loaded once.
func ListWithError() ([]Variety, error) {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}

	return append([]Variety(nil), config.Values...), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
func List() []Variety {
	values, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.Error("error.loading:file", map[string]any{
				"Subject": "variety",
				"File":    fileConfigPath,
				"Error":   err,
			})
		})
		return nil
	}

	return values
}

// Validate reports whether r is a configured value.
func (r Variety) Validate() bool {
	return slices.Contains(List(), r)
}

// Default returns the configured default variety, or registers "card" when unset or invalid.
func Default() Variety {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return MakeDefault("card")
	}

	if config.Default == "" || !config.Default.Validate() {
		return MakeDefault("card")
	}

	return config.Default
}

// MakeDefault adds defaultValue to the cached values when missing and returns it.
func MakeDefault(defaultValue string) Variety {
	_, _ = ListWithError()
	cachedConfig.Update(func(config *yamlConfig) {
		config.addValue(Variety(defaultValue))
	})
	return Variety(defaultValue)
}

// AssertWithDefault returns provided when valid, otherwise Default.
func AssertWithDefault(provided string) Variety {
	if Assert(provided) {
		return Variety(provided)
	}

	return Default()
}

// Assert reports whether provided is a configured value.
func Assert(provided string) bool {
	return Variety(provided).Validate()
}

// String returns r as a string.
func (r Variety) String() string {
	return string(r)
}
