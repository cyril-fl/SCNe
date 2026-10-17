package locations

import (
	"errors"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"solopg/config"
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

var cachedConfig []yamlConfig
var lockedNames map[string]struct{}

func loadFromSource() error {
	lockedNames = make(map[string]struct{})
	files, err := yaml.GetFilesFromSource(folderConfigPath, true)
	if err != nil {
		return logs.Error("error.locations.load_folder", map[string]any{"Folder": folderConfigPath, "Error": err})
	}

	if errs := handleLoadFromFiles(files); len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func handleLoadFromFiles(files []string) []error {
	errs := []error{}
	for _, file := range files {
		if err := loadFromFile(file); err != nil {
			errs = append(errs, logs.CeaseError("error.loading:file", map[string]any{
				"Subject": "location",
				"File":    file,
				"Error":   err,
			}))
			continue
		}
	}

	return errs
}

func loadFromFile(fileAddress string) error {
	params, err := yaml.LoadFromFile[yamlConfig](fileAddress)

	if err != nil {

		return logs.Error("error.loading:file", map[string]any{
			"Subject": "location",
			"File":    fileAddress,
			"Error":   err,
		})
	}

	/*
		LOW Verrifier ou "MakeDefault" pourrait être utile
		ex: Si la data vien d'une carte loader, c'est valider

		Exemple a suivre pour les autres variete et autres pourquoi ? si ca vien d'une carte loader, c'est valider !
	*/
	variety.MakeDefault(params.Variety)

	cachedConfig = append(cachedConfig, *params)
	lockedNames[params.Name] = struct{}{}

	return nil
}

// - Collection - //
func List() []Location {
	if cachedConfig == nil {
		err := loadFromSource()
		if err != nil {
			logs.SilentWarning("error.loading", map[string]any{
				"Subject": "location",
				"Error":   err,
			})
			return nil
		}
	}

	return makeSet(cachedConfig)
}

// IsLocked reports whether a location comes from the initial template cache.
func IsLocked(name string) bool {
	if cachedConfig == nil {
		List()
	}
	_, ok := lockedNames[name]
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
