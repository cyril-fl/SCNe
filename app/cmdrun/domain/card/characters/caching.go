package characters

import (
	"errors"
	"fmt"
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
	"solopg/app/shared/types/primitive"
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

type cacheValue struct {
	Configs     []yamlConfig
	LockedNames map[string]struct{}
}

var cachedConfig = primitive.NewCache[cacheValue]()

func loadFromSource() (cacheValue, error) {
	files, err := yaml.GetFilesFromSource(folderConfigPath, true)
	if err != nil {
		return cacheValue{}, err
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
	params, err := yaml.LoadFromFile[yamlConfig](fileAddress)

	if err != nil {
		return yamlConfig{}, fmt.Errorf("load %s: %w", fileAddress, err)
	}

	/*
		LOW Verrifier ou "MakeDefault" pourrait être utile
		ex: Si la data vien d'une carte loader, c'est valider
	*/
	variety.MakeDefault(params.Variety)

	return *params, nil
}

// - Collection - //
func ListWithError() ([]Character, error) {
	config, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return nil, err
	}

	return makeSet(config.Configs), nil
}

func List() []Character {
	characters, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.SilentWarning("error.loading", map[string]any{
				"Subject": "character",
				"Folder":  folderConfigPath,
				"Error":   err,
			})
		})
		return nil
	}

	return characters
}

// IsLocked reports whether a character comes from the initial template cache.
func IsLocked(name string) bool {
	config, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return false
	}
	_, ok := config.LockedNames[name]
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
