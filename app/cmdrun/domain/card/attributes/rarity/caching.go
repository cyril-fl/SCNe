package rarity

import (
	"slices"
	"solopg/app/shared/services/yaml"
	"solopg/config"
)

// - Configuration & caching - //
type yamlConfig struct {
	Values  []Rarity `yaml:"values"`
	Default Rarity   `yaml:"default"`
}

var fileConfigPath = config.Current.Documents.Files.Rarity

var cachedConfig yamlConfig

func loadFromFile() error {
	params, err := yaml.LoadFromFile[yamlConfig](fileConfigPath)
	if err != nil {
		return err
	}

	cachedConfig = *params

	return nil
}

func (c *yamlConfig) addValue(value Rarity) {
	if !slices.Contains(c.Values, value) {
		c.Values = append(c.Values, value)
	}
}

// - Rarity - //
type Rarity string

func List() []Rarity {
	if len(cachedConfig.Values) == 0 {
		err := loadFromFile()
		if err != nil {
			return nil
		}
	}

	return cachedConfig.Values
}

func Default() Rarity {
	if cachedConfig.Default == "" {
		if err := loadFromFile(); err != nil || !cachedConfig.Default.Validate() {
			return MakeDefault("")
		}
	}

	return cachedConfig.Default
}

func (r Rarity) Validate() bool {
	return slices.Contains(List(), r)
}

func MakeDefault(defaultValue string) Rarity {
	List()
	cachedConfig.addValue(Rarity(defaultValue))
	return Rarity(defaultValue)
}

func AssertWithDefault(provided string) Rarity {
	if Assert(provided) {
		return Rarity(provided)
	}

	return Default()
}

func Assert(provided string) bool {
	return Rarity(provided).Validate()
}

func (r Rarity) String() string {
	return string(r)
}
