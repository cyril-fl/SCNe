package rank

import (
	"slices"
	"solopg/app/shared/services/logs"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Values []Rank `yaml:"values"`
}

var fileConfigPath = config.GetConfig().Documents.Files.Ranks

var cachedConfig = plumprimitive.NewCache[yamlConfig]()

func loadFromFile() (yamlConfig, error) {
	loader := plummillload.New(fileConfigPath)
	params := loader.LoadFromFile[yamlConfig]()
	if loader.HasErr() {
		return yamlConfig{}, logs.RawError("%w", loader.GetErr())
	}

	return *params, nil
}

func (c *yamlConfig) addValue(value Rank) {
	if !slices.Contains(c.Values, value) {
		c.Values = append(c.Values, value)
	}
}

// - Rank - //
type Rank string

func ListWithError() ([]Rank, error) {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}

	return append([]Rank(nil), config.Values...), nil
}

func List() []Rank {
	values, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.Error("error.loading:file", map[string]any{
				"Subject": "rank",
				"File":    fileConfigPath,
				"Error":   err,
			})
		})
		return nil
	}

	return values
}

func (r Rank) Validate() bool {
	return slices.Contains(List(), r)
}

func Default() Rank {
	list := List()

	if len(list) == 0 {
		return MakeDefault("side")
	}

	return list[0]
}

func MakeDefault(defaultValue string) Rank {
	_, _ = ListWithError()
	cachedConfig.Update(func(config *yamlConfig) {
		config.addValue(Rank(defaultValue))
	})
	return Rank(defaultValue)
}

func AssertWithDefault(provided string) Rank {
	if Assert(provided) {
		return Rank(provided)
	}

	return Default()
}

func Assert(provided string) bool {
	return Rank(provided).Validate()
}

func (r Rank) String() string {
	return string(r)
}
