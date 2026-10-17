package dice

import (
	"math/rand"

	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"solopg/app/shared/types/primitive"
	"solopg/config"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name  string `yaml:"name"`
	Sides int    `yaml:"sides"`
}

var fileConfigPath = config.Current.Documents.Files.Dice

var cachedConfig = primitive.NewCache[[]yamlConfig]()

func loadFromFile() ([]yamlConfig, error) {
	params, err := yaml.LoadListFromFile[yamlConfig](fileConfigPath)
	if err != nil {
		return nil, err
	}
	return params, nil
}

// - Dice - //
type Dice struct {
	name  string
	sides int
}

func new(config yamlConfig) Dice {
	return Dice{
		name:  config.Name,
		sides: config.Sides,
	}
}

func (d Dice) GetName() string {
	if d.name != "" {
		return d.name
	}
	return i19n.Localize("dice.default", map[string]any{
		"sides": d.sides,
	})
}

func (d Dice) GetSides() int {
	return d.sides
}

func (d Dice) Roll() int {
	return Roll(d.sides)
}

// - Collection - //
func ListWithError() ([]Dice, error) {
	configs, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}
	return makeSet(configs), nil
}

func List() []Dice {
	dice, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.SilentWarning("error.loading:file", map[string]any{
				"Subject": "dice",
				"File":    fileConfigPath,
				"Error":   err,
			})
		})
		return nil
	}
	return dice
}

func makeSet(configs []yamlConfig) []Dice {
	result := make([]Dice, 0, len(configs))

	for _, config := range configs {
		result = append(result, new(config))
	}

	return result
}

// - Helpers - //
func Roll(sides int) int {
	if sides <= 0 {
		return 0
	}
	return 1 + (rand.Intn(sides))
}
