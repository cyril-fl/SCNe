package locations

import (
	"errors"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/shared/services/logs"
	"solopg/config"

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

var folderConfigPath = config.Current.Documents.Folders.Locations

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

	configs := make([]yamlConfig, 0, len(files))
	locked := make(map[string]struct{}, len(files))
	var errs []error
	for _, file := range files {
		config, err := loadFromFile(file)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		configs = append(configs, config)
		locked[config.Name] = struct{}{}
	}

	if err := errors.Join(errs...); err != nil {
		return cacheValue{}, err
	}

	return cacheValue{Configs: configs, LockedNames: locked}, nil
}

func loadFromFile(fileAddress string) (yamlConfig, error) {
	loader := plummillload.New(fileAddress)
	params := loader.LoadFromFile[yamlConfig]()
	if loader.HasErr() {
		return yamlConfig{}, logs.RawError("load %s: %w", fileAddress, loader.GetErr())
	}

	/*
		LOW Verrifier ou "MakeDefault" pourrait être utile
		ex: Si la data vien d'une carte loader, c'est valider

		Exemple a suivre pour les autres variete et autres pourquoi ? si ca vien d'une carte loader, c'est valider !
	*/
	variety.MakeDefault(params.Variety)

	return *params, nil
}

// - Collection - //
func ListWithError() ([]Location, error) {
	config, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return nil, err
	}

	return makeSet(config.Configs), nil
}

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
		Rarity:      rarity.MakeDefault(config.Rarity),
		Variety:     variety.MakeDefault(config.Variety),
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
