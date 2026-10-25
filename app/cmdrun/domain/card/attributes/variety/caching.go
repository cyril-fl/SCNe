package variety

import (
	"slices"
	"solopg/app/shared/services/logs"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Values []Variety `yaml:"values"`
}

var fileConfigPath = config.GetConfig().Documents.Files.Variety

var cachedConfig = plumprimitive.NewCache[yamlConfig]()

func loadFromFile() (yamlConfig, error) {
	loader := plummillload.New(fileConfigPath)
	params := loader.LoadFromFile[yamlConfig]()
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

// Default returns the first configured variety, or registers "card" when none.
func Default() Variety {
	list := List()

	if len(list) == 0 {
		return MakeDefault("card")
	}

	return list[0]
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
