package variety

import (
	"slices"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"solopg/app/shared/types/primitive"
	"solopg/config"
)

// - Configuration & caching - //
type yamlConfig struct {
	Values []Variety `yaml:"values"`
}

var fileConfigPath = config.Current.Documents.Files.Variety

var cachedConfig = primitive.NewCache[yamlConfig]()

func loadFromFile() (yamlConfig, error) {
	params, err := yaml.LoadFromFile[yamlConfig](fileConfigPath)
	if err != nil {
		return yamlConfig{}, err
	}

	return *params, nil
}

func (c *yamlConfig) addValue(value Variety) {
	if !slices.Contains(c.Values, value) {
		c.Values = append(c.Values, value)
	}
}

// - Variety - //
type Variety string

func ListWithError() ([]Variety, error) {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}

	return append([]Variety(nil), config.Values...), nil
}

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
	_, _ = ListWithError()
	cachedConfig.Update(func(config *yamlConfig) {
		config.addValue(Variety(defaultValue))
	})
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
