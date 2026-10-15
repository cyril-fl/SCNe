package races

import (
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"solopg/app/shared/types/primitive"
	"solopg/config"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Playable    bool             `yaml:"playable"`
	Bonus       []stats.Modifier `yaml:"bonus"`
}

var fileConfigPath = config.Current.Documents.Files.Races

var cachedConfig []yamlConfig
var lockedNames map[string]struct{}

func loadFromFile() error {
	params, err := yaml.LoadListFromFile[yamlConfig](fileConfigPath)
	if err != nil {
		return logs.Error("error.loading:file", map[string]any{
			"Subject": "race",
			"File":    fileConfigPath,
			"Error":   err,
		})
	}

	cachedConfig = params
	lockedNames = make(map[string]struct{}, len(params))
	for _, value := range params {
		lockedNames[value.Name] = struct{}{}
	}

	return nil
}

func Has(value Template) bool {
	new := New(value)

	for _, race := range List() {
		if race.GetHash() == new.GetHash() {
			return true
		}
	}
	return false
}

func AddInConfig(value Template) error {
	if Has(value) {
		return logs.MildError("error.invalid:already_exists", map[string]any{
			"Subject": "race",
			"Value":   value.Name,
		})
	}

	new := yamlConfig{
		Name:        value.Name,
		Description: value.Description,
		Playable:    value.Playable,
		Bonus:       value.Bonus,
	}

	cachedConfig = append(cachedConfig, new)

	return nil
}

// - Collection - //
func List() []Race {
	if cachedConfig == nil {
		if err := loadFromFile(); err != nil {
			logs.SilentWarning("error.loading", map[string]any{
				"Subject": "race",
				"Error":   err,
			})
			return nil
		}
	}
	return makeSet(cachedConfig)
}

func makeSet(configs []yamlConfig) []Race {
	result := make([]Race, 0, len(configs))

	for _, config := range configs {
		result = append(result, newFromYaml(config))
	}

	return result
}

func newFromYaml(config yamlConfig) Race {
	return Race{
		Descriptible: primitive.Describe(primitive.TemplateDescriptible{
			Name:        config.Name,
			Description: config.Description,
		}),
		Playable:    config.Playable,
		Bonus:       config.Bonus,
	}
}

func ListNames() []string {
	var list []string
	for _, i := range List() {
		list = append(list, i.GetName())
	}
	return list
}

// IsLocked reports whether a race comes from the initial template cache.
// Races added later with AddInConfig are intentionally not locked.
func IsLocked(name string) bool {
	if cachedConfig == nil {
		List()
	}
	_, ok := lockedNames[name]
	return ok
}

func FindByName(name string) *Race {
	for _, i := range List() {
		if i.GetName() == name {
			return &i
		}
	}
	return nil
}
