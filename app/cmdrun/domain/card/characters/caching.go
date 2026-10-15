package characters

import (
	"errors"
	"os"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/card/characters/wallet"
	"solopg/app/cmdrun/domain/card/objects"
	"solopg/app/cmdrun/domain/card/objects/equipment"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"solopg/config"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Rarity      string              `yaml:"rarity"`
	Variety     string              `yaml:"variety"`
	Class       string              `yaml:"class"`
	Race        string              `yaml:"race"`
	Stats       stats.Stats         `yaml:"stats"`
	Equipment   equipment.Equipment `yaml:"equipment"`
	Inventory   []objects.Object    `yaml:"inventory"`
	Wallet      wallet.Wallet       `yaml:"wallet"`
}

var folderConfigPath = config.Current.Documents.Folders.Characters

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
				"Subject": "character",
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
		return logs.Error("error.loading", map[string]any{
			"Subject": "location",
			"Error":   err,
		})
	}

	/*
		LOW Verrifier ou "MakeDefault" pourrait être utile
		ex: Si la data vien d'une carte loader, c'est valider
	*/
	variety.MakeDefault(params.Variety)

	cachedConfig = append(cachedConfig, *params)
	lockedNames[params.Name] = struct{}{}

	return nil
}

// - Collection - //
func List() []Character {
	if cachedConfig == nil {
		if err := loadFromSource(); err != nil {
			logs.SilentWarning("error.loading", map[string]any{
				"Subject": ("dice"),
				"Error":   err,
			})
			return nil
		}
	}

	return makeSet(cachedConfig)
}

// IsLocked reports whether a character comes from the initial template cache.
func IsLocked(name string) bool {
	if cachedConfig == nil {
		List()
	}
	_, ok := lockedNames[name]
	return ok
}

func makeSet(configs []yamlConfig) []Character {
	result := make([]Character, 0, len(configs))

	for _, config := range configs {
		if location := newFromYaml(config); location != nil {
			result = append(result, *location)
		}
	}

	return result
}

func newFromYaml(config yamlConfig) *Character {
	class, err := classes.InsertIfIsnt(config.Class)
	if err != nil {
		cwd, cwderr := os.Getwd()
		logs.SilentWarning("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, err),
		})
		return nil
	}

	newLocation, err := New(Template{
		Name:        config.Name,
		Description: config.Description,
		Class:       class.GetName(),
		Race:        config.Race,
		Rarity:      rarity.MakeDefault(config.Rarity),
		Variety:     variety.MakeDefault(config.Variety),
		Stats:       config.Stats,
		Wallet:      config.Wallet,
		Equipment:   config.Equipment,
		Inventory:   config.Inventory,
	})

	if err != nil {
		cwd, cwderr := os.Getwd()
		logs.SilentWarning("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, err),
		})
		return nil
	}

	return newLocation
}
