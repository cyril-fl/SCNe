package rarity

import (
	"slices"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"solopg/app/shared/types/primitive"
	"solopg/config"
)

// - Configuration & caching - //
type yamlConfig struct {
	Values  []Rarity `yaml:"values"`
	Default Rarity   `yaml:"default"`
}

var fileConfigPath = config.Current.Documents.Files.Rarity

var cachedConfig = primitive.NewCache[yamlConfig]()

func loadFromFile() (yamlConfig, error) {
	params, err := yaml.LoadFromFile[yamlConfig](fileConfigPath)
	if err != nil {
		return yamlConfig{}, err
	}

	return *params, nil
}

func (c *yamlConfig) addValue(value Rarity) {
	if !slices.Contains(c.Values, value) {
		c.Values = append(c.Values, value)
	}
}

// - Rarity - //
type Rarity string

func ListWithError() ([]Rarity, error) {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}

	return append([]Rarity(nil), config.Values...), nil
}

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

func (r Rarity) Validate() bool {
	return slices.Contains(List(), r)
}

func MakeDefault(defaultValue string) Rarity {
	_, _ = ListWithError()
	cachedConfig.Update(func(config *yamlConfig) {
		config.addValue(Rarity(defaultValue))
	})
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
