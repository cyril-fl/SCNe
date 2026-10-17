package equipment

import (
	"errors"
	"solopg/app/cmdrun/domain/card/attributes/objectcategory"
	"solopg/app/cmdrun/types/potency"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/slot"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name        string                  `yaml:"name"`
	Description string                  `yaml:"description"`
	Set         string                  `yaml:"set"`
	Rarity      rarity.Rarity           `yaml:"rarity"`
	Variety     variety.Variety         `yaml:"variety"`
	Category    objectcategory.Category `yaml:"category"`
	Potency     potency.Value           `yaml:"potency"`
	Effects     []stats.Effect          `yaml:"effects"`
	Pod         int                     `yaml:"pod"`
	Slot        slot.Slot               `yaml:"destinedslot"`
}

var folderConfigPath = "data/template/armor_sets"

var cachedConfig []yamlConfig

func loadFromSource() error {
	files, err := yaml.GetFilesFromSource(folderConfigPath, true)
	if err != nil {

		return logs.Error("error.loading:folder", map[string]any{
			"Subject": "gear",
			"Folder":  folderConfigPath,
			"Error":   err,
		})
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

			errs = append(errs, logs.Error("error.loading:file", map[string]any{
				"Subject": "gear",
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
		return err
	}

	cachedConfig = append(cachedConfig, *params)

	return nil
}

// - Collection - //
func List() []Gear {
	if cachedConfig == nil {
		if err := loadFromSource(); err != nil {
			logs.SilentWarning("error.loading", map[string]any{
				"Subject": "gear",
				"Error":   err,
			})
			return nil
		}
	}

	return makeSet(cachedConfig)
}

func makeSet(configs []yamlConfig) []Gear {
	result := make([]Gear, 0, len(configs))

	for _, config := range configs {
		if location := newFromYaml(config); location != nil {
			result = append(result, *location)
		}
	}

	return result
}

func newFromYaml(config yamlConfig) *Gear {
	newGear, err := NewGear(Template{
		Name:        config.Name,
		Description: config.Description,
		Rarity:      config.Rarity,
		Variety:     config.Variety,
		Category:    config.Category,
		Potency:     config.Potency,
		Effects:     config.Effects,
		Pod:         config.Pod,
		Slot:        config.Slot,
	})

	if err != nil {
		logs.SilentWarning("error.invalid:new", map[string]any{
			"Subject": "gear",
			"Error":   err,
		})
		return nil
	}

	return newGear
}

func FindEquipementByName(name string) []Gear {
	if cachedConfig == nil {
		if err := loadFromSource(); err != nil {
			logs.SilentWarning("error.loading", map[string]any{
				"Subject": "gear",
				"Error":   err,
			})
			return nil
		}
	}

	var serchResults []yamlConfig
	for _, gear := range cachedConfig {
		if gear.Set == name {
			serchResults = append(serchResults, gear)
		}
	}

	return makeSet(serchResults)
}
