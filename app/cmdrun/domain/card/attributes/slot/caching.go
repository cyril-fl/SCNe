package slot

import (
	"slices"
	"solopg/app/shared/services/yaml"
	"solopg/config"
)

// - Configuration & caching - //
type yamlConfig struct {
	Values []Slot `yaml:"values"`
}

var fileConfigPath = config.Current.Documents.Files.Slots

var cachedConfig yamlConfig

func loadFromFile() error {
	params, err := yaml.LoadFromFile[yamlConfig](fileConfigPath)
	if err != nil {
		return err
	}

	cachedConfig = *params

	return nil
}

func (c *yamlConfig) addValue(value Slot) {
	if !slices.Contains(c.Values, value) {
		c.Values = append(c.Values, value)
	}
}

// Slot
type Slot string

func List() []Slot {
	if len(cachedConfig.Values) == 0 {
		if err := loadFromFile(); err != nil {
			return nil
		}
	}

	return cachedConfig.Values
}

func (r Slot) Validate() bool {
	return slices.Contains(List(), r)
}

func Assert(provided string) bool {
	return Slot(provided).Validate()
}

func (r Slot) String() string {
	return string(r)
}
