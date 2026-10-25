package locations

import (
	"errors"
	"scne/app/cmdrun/domain/card/attributes/rarity"
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/attributes/variety"
	"scne/app/shared/services/logs"
	"scne/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillentries "github.com/cyril-fl/plum/services/factory/millentries"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Rarity      string         `yaml:"rarity"`
	Variety     string         `yaml:"variety"`
	Effects     []stats.Effect `yaml:"effects"`
}

var folderConfigPath = config.GetConfig().Documents.Folders.Locations

type cacheValue struct {
	Configs     []yamlConfig
	LockedNames map[string]struct{}
}

var cachedConfig = plumprimitive.NewCache[cacheValue]()

func loadFromSource() (cacheValue, error) {
	entries := plummillentries.New()
	files := entries.GetFilesFromSource(folderConfigPath, true)
	if entries.HasErr() {
		return cacheValue{}, logs.RawError("%w", entries.GetErr())
	}

	configs, errs := loadFromFiles(files)

	if len(errs) > 0 {
		return cacheValue{}, errors.Join(errs...)
	}

	locked := lockNames(configs)

	return cacheValue{Configs: configs, LockedNames: locked}, nil
}

func lockNames(configs []yamlConfig) map[string]struct{} {
	locked := make(map[string]struct{}, len(configs))

	for _, config := range configs {
		locked[config.Name] = struct{}{}
	}

	return locked
}

func loadFromFiles(files []string) ([]yamlConfig, []error) {
	configs := make([]yamlConfig, 0, len(files))
	errs := make([]error, 0, len(files))

	for _, file := range files {
		config, err := loadFromFile(file)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		configs = append(configs, config)
	}

	return configs, errs
}

func loadFromFile(fileAddress string) (yamlConfig, error) {
	loader := plummillload.New(fileAddress)
	params := loader.LoadFromFile[yamlConfig]()
	if loader.HasErr() {
		return yamlConfig{}, logs.RawError("load %s: %w", fileAddress, loader.GetErr())
	}

	// Values from a loaded card file are trusted: register them as valid.
	if params.Rarity != "" {
		rarity.MakeDefault(params.Rarity)
	}
	if params.Variety != "" {
		variety.MakeDefault(params.Variety)
	}

	return *params, nil
}

// - Collection - //

// ListWithError returns the locations loaded from the locations folder, once.
// Invalid cards are logged and skipped.
func ListWithError() ([]Location, error) {
	config, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return nil, err
	}

	return makeSet(config.Configs), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
func List() []Location {
	locations, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.SilentWarning("error.loading", map[string]any{
				"Subject": "location",
				"Folder":  folderConfigPath,
				"Error":   err,
			})
		})
		return nil
	}

	return locations
}

// IsLocked reports whether a location comes from the initial template cache.
func IsLocked(name string) bool {
	config, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return false
	}
	_, ok := config.LockedNames[name]
	return ok
}

func makeSet(configs []yamlConfig) []Location {
	result := make([]Location, 0, len(configs))

	for _, config := range configs {
		if location := newFromYaml(config); location != nil {
			result = append(result, *location)
		}
	}

	return result
}

func newFromYaml(config yamlConfig) *Location {
	newLocation, err := New(Template{
		Name:        config.Name,
		Description: config.Description,
		Rarity:      rarity.Rarity(config.Rarity),
		Variety:     variety.Variety(config.Variety),
		Effects:     config.Effects,
	})

	if err != nil {
		logs.SilentWarning("error.invalid:new", map[string]any{
			"Subject": "location",
			"Error":   err,
		})
		return nil
	}

	return newLocation
}
