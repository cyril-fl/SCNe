package races

import (
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/shared/services/logs"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Playable    bool             `yaml:"playable"`
	Bonus       []stats.Modifier `yaml:"bonus"`
}

var fileConfigPath = config.GetConfig().Documents.Files.Races

type cacheValue struct {
	Configs     []yamlConfig
	LockedNames map[string]struct{}
}

var cachedConfig = plumprimitive.NewCache[cacheValue]()

func loadFromFile() ([]yamlConfig, error) {
	loader := plummillload.New(fileConfigPath)
	params := loader.LoadListFromFile[yamlConfig]()
	if loader.HasErr() {
		return nil, logs.RawError("%w", loader.GetErr())
	}

	return params, nil
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
	if _, err := ListWithError(); err != nil {
		reportLoadError()
		return err
	}

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

	cachedConfig.Update(func(config *cacheValue) {
		config.Configs = append(config.Configs, new)
	})

	return nil
}

// - Collection - //
func ListWithError() ([]Race, error) {
	config, err := cachedConfig.Load(func() (cacheValue, error) {
		configs, err := loadFromFile()
		if err != nil {
			return cacheValue{}, err
		}

		locked := make(map[string]struct{}, len(configs))
		for _, value := range configs {
			locked[value.Name] = struct{}{}
		}

		return cacheValue{Configs: configs, LockedNames: locked}, nil
	})
	if err != nil {
		return nil, err
	}

	return makeSet(config.Configs), nil
}

func List() []Race {
	races, err := ListWithError()
	if err != nil {
		reportLoadError()
		return nil
	}

	return races
}

func reportLoadError() {
	cachedConfig.ReportError(func(err error) {
		logs.SilentWarning("error.loading", map[string]any{
			"Subject": "race",
			"File":    fileConfigPath,
			"Error":   err,
		})
	})
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
		Descriptible: plumprimitive.Describe(plumprimitive.TemplateDescriptible{
			Name:        config.Name,
			Description: config.Description,
		}),
		Playable: config.Playable,
		Bonus:    config.Bonus,
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
	config, err := cachedConfig.Load(func() (cacheValue, error) {
		configs, err := loadFromFile()
		if err != nil {
			return cacheValue{}, err
		}

		locked := make(map[string]struct{}, len(configs))
		for _, value := range configs {
			locked[value.Name] = struct{}{}
		}

		return cacheValue{Configs: configs, LockedNames: locked}, nil
	})
	if err != nil {
		reportLoadError()
		return false
	}
	_, ok := config.LockedNames[name]
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
