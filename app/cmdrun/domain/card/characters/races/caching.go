package races

import (
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/shared/services/logs"
	"scne/config"
	"slices"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	ID          string           `yaml:"id"`
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Playable    bool             `yaml:"playable"`
	Bonus       []stats.Modifier `yaml:"bonus"`
}

var fileConfigPath = config.GetConfig().Documents.Files.Races

type cacheValue struct {
	Configs   []yamlConfig
	LockedIDs map[string]struct{}
}

var cachedConfig = plumprimitive.NewCache[cacheValue]()

func loadFromFile() ([]yamlConfig, error) {
	loader := plummillload.New(fileConfigPath)
	params := loader.LoadListFromFile[yamlConfig]()
	if loader.HasErr() {
		return nil, logs.RawError("%w", loader.GetErr())
	}

	// Template races are never renamed, so their name is a stable ID.
	for i := range params {
		if params[i].ID == "" {
			params[i].ID = params[i].Name
		}
	}

	return params, nil
}

// Has reports whether a race with value.ID is configured.
func Has(value Template) bool {
	return FindByID(value.ID) != nil
}

// AddInConfig adds value to the cached races; it fails when its ID exists.
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
		ID:          value.ID,
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

// ReplaceInConfig swaps the race with this id for value, or adds value when
// id is not in the config.
func ReplaceInConfig(id string, value Template) error {
	if _, err := ListWithError(); err != nil {
		reportLoadError()
		return err
	}

	replacement := yamlConfig{
		ID:          value.ID,
		Name:        value.Name,
		Description: value.Description,
		Playable:    value.Playable,
		Bonus:       value.Bonus,
	}

	cachedConfig.Update(func(config *cacheValue) {
		index := slices.IndexFunc(config.Configs, func(c yamlConfig) bool { return c.ID == id })
		if index < 0 {
			config.Configs = append(config.Configs, replacement)
			return
		}
		config.Configs[index] = replacement
	})

	return nil
}

// - Collection - //

// ListWithError returns the configured races; the file is loaded once.
func ListWithError() ([]Race, error) {
	config, err := cachedConfig.Load(func() (cacheValue, error) {
		configs, err := loadFromFile()
		if err != nil {
			return cacheValue{}, err
		}

		locked := make(map[string]struct{}, len(configs))
		for _, value := range configs {
			locked[value.ID] = struct{}{}
		}

		return cacheValue{Configs: configs, LockedIDs: locked}, nil
	})
	if err != nil {
		return nil, err
	}

	return makeSet(config.Configs), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
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
		ID:       config.ID,
		Playable: config.Playable,
		Bonus:    config.Bonus,
	}
}

// IsLocked reports whether a race comes from the initial template cache.
// Races added later with AddInConfig are intentionally not locked.
func IsLocked(id string) bool {
	config, err := cachedConfig.Load(func() (cacheValue, error) {
		configs, err := loadFromFile()
		if err != nil {
			return cacheValue{}, err
		}

		locked := make(map[string]struct{}, len(configs))
		for _, value := range configs {
			locked[value.ID] = struct{}{}
		}

		return cacheValue{Configs: configs, LockedIDs: locked}, nil
	})
	if err != nil {
		reportLoadError()
		return false
	}
	_, ok := config.LockedIDs[id]
	return ok
}

// FindByID returns the race with id, or nil.
func FindByID(id string) *Race {
	for _, i := range List() {
		if i.GetID() == id {
			return &i
		}
	}
	return nil
}

// NameOf returns the name of the race with this id, or id when it is unknown.
func NameOf(id string) string {
	if race := FindByID(id); race != nil {
		return race.GetName()
	}
	return id
}
