package dice

import (
	"math/rand"

	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"
	"scne/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name  string `yaml:"name"`
	Sides int    `yaml:"sides"`
}

var fileConfigPath = config.GetConfig().Documents.Files.Dice

var cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

func loadFromFile() ([]yamlConfig, error) {
	loader := plummillload.New(fileConfigPath)
	params := loader.LoadListFromFile[yamlConfig]()
	if loader.HasErr() {
		return nil, logs.RawError("%w", loader.GetErr())
	}
	return params, nil
}

// - Dice - //

// Dice is a configured die.
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

// GetName returns the configured name, or a localized default from the sides.
func (d Dice) GetName() string {
	if d.name != "" {
		return d.name
	}
	return i19n.Localize("dice.default", map[string]any{
		"sides": d.sides,
	})
}

// GetSides returns the number of sides.
func (d Dice) GetSides() int {
	return d.sides
}

// Roll rolls the die.
func (d Dice) Roll() int {
	return Roll(d.sides)
}

// - Collection - //

// ListWithError returns the configured dice; the file is loaded once.
func ListWithError() ([]Dice, error) {
	configs, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}
	return makeSet(configs), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
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

// Roll returns a random value in [1, sides], or 0 when sides is not positive.
func Roll(sides int) int {
	if sides <= 0 {
		return 0
	}
	return 1 + (rand.Intn(sides))
}
