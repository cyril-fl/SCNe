package equipment

import (
	"errors"
	"fmt"

	"solopg/app/cmdrun/domain/card/attributes/objectcategory"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/slot"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/cmdrun/types/potency"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"solopg/app/shared/types/primitive"
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

var cachedConfig = primitive.NewCache[[]yamlConfig]()

func loadFromSource() ([]yamlConfig, error) {
	files, err := yaml.GetFilesFromSource(folderConfigPath, true)
	if err != nil {
		return nil, fmt.Errorf("load gear folder %s: %w", folderConfigPath, err)
	}

	configs := make([]yamlConfig, 0, len(files))
	errs := make([]error, 0)
	for _, file := range files {
		config, err := loadFromFile(file)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		configs = append(configs, config)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return configs, nil
}

func loadFromFile(fileAddress string) (yamlConfig, error) {
	params, err := yaml.LoadFromFile[yamlConfig](fileAddress)
	if err != nil {
		return yamlConfig{}, fmt.Errorf("load gear file %s: %w", fileAddress, err)
	}
	return *params, nil
}

// - Collection - //
func ListWithError() ([]Gear, error) {
	configs, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return nil, err
	}
	return makeSet(configs), nil
}

func List() []Gear {
	gears, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.SilentWarning("error.loading", map[string]any{
				"Subject": "gear",
				"Error":   err,
			})
		})
		return nil
	}
	return gears
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
	configs, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return nil
	}

	searchResults := make([]yamlConfig, 0)
	for _, gear := range configs {
		if gear.Set == name {
			searchResults = append(searchResults, gear)
		}
	}

	return makeSet(searchResults)
}
