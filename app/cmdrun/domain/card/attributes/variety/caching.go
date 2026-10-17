package variety

import (
	"slices"
	"solopg/app/shared/services/yaml"
	"solopg/config"
)

// - Configuration & caching - //
type yamlConfig struct {
	Values []Variety `yaml:"values"`
}

var fileConfigPath = config.Current.Documents.Files.Variety

var cachedConfig yamlConfig

func loadFromFile() error {
	params, err := yaml.LoadFromFile[yamlConfig](fileConfigPath)
	if err != nil {
		return err
	}

	cachedConfig = *params

	return nil
}

func (c *yamlConfig) addValue(value Variety) {
	if !slices.Contains(c.Values, value) {
		c.Values = append(c.Values, value)
	}
}

// - Variety - //
type Variety string

func List() []Variety {
	if len(cachedConfig.Values) == 0 {
		if err := loadFromFile(); err != nil {
			return nil
		}
	}

	return cachedConfig.Values
}

func (r Variety) Validate() bool {
	return slices.Contains(List(), r)
}

func Default() Variety {
	list := List()

	if len(list) == 0 {
		return MakeDefault("card")
	}

	return list[0]
}

func MakeDefault(defaultValue string) Variety {
	List()
	cachedConfig.addValue(Variety(defaultValue))
	return Variety(defaultValue)
}

func AssertWithDefault(provided string) Variety {
	if Assert(provided) {
		return Variety(provided)
	}

	return Default()
}

func Assert(provided string) bool {
	return Variety(provided).Validate()
}

func (r Variety) String() string {
	return string(r)
}
