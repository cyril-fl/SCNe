package slot

import (
	"scne/app/shared/services/logs"
	"scne/config"
	"slices"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Values []Slot `yaml:"values"`
}

var fileConfigPath = config.GetConfig().Documents.Files.Slots

var cachedConfig = plumprimitive.NewCache[yamlConfig]()

func loadFromFile() (yamlConfig, error) {
	loader := plummillload.New(fileConfigPath)
	params := loader.LoadFromFile[yamlConfig]()
	if loader.HasErr() {
		return yamlConfig{}, logs.RawError("%w", loader.GetErr())
	}

	return *params, nil
}

func (c *yamlConfig) addValue(value Slot) {
	if !slices.Contains(c.Values, value) {
		c.Values = append(c.Values, value)
	}
}

// Slot

// Slot is a value from the slot config file.
type Slot string

// ListWithError returns a copy of the configured values; the file is loaded once.
func ListWithError() ([]Slot, error) {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}

	return append([]Slot(nil), config.Values...), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
func List() []Slot {
	values, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.Error("error.loading:file", map[string]any{
				"Subject": "slot",
				"File":    fileConfigPath,
				"Error":   err,
			})
		})
		return nil
	}

	return values
}

// Validate reports whether r is a configured value.
func (r Slot) Validate() bool {
	return slices.Contains(List(), r)
}

// Assert reports whether provided is a configured value.
func Assert(provided string) bool {
	return Slot(provided).Validate()
}

// String returns r as a string.
func (r Slot) String() string {
	return string(r)
}
