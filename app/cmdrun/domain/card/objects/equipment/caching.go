package equipment

import (
	"errors"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillentries "github.com/cyril-fl/plum/services/factory/millentries"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
	"solopg/app/cmdrun/domain/card/attributes/objectcategory"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/slot"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/cmdrun/types/potency"
	"solopg/app/shared/services/logs"
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

var cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

func loadFromSource() ([]yamlConfig, error) {
	entries := plummillentries.New()
	files := entries.GetFilesFromSource(folderConfigPath, true)
	if entries.HasErr() {
		return nil, logs.RawError("load gear folder %s: %w", folderConfigPath, entries.GetErr())
	}

	configs, errs := loadFromFiles(files)

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return configs, nil
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
		return yamlConfig{}, logs.RawError("load gear file %s: %w", fileAddress, loader.GetErr())
	}

	// Values from a loaded card file are trusted: register them as valid.
	if params.Rarity != "" {
		rarity.MakeDefault(string(params.Rarity))
	}
	if params.Variety != "" {
		variety.MakeDefault(string(params.Variety))
	}
	if params.Category != "" {
		objectcategory.MakeDefault(string(params.Category))
	}

	return *params, nil
}

// - Collection - //

// ListWithError returns the gears loaded from the armor sets folder, once.
// Invalid gears are logged and skipped.
func ListWithError() ([]Gear, error) {
	configs, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return nil, err
	}
	return makeSet(configs), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
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

// FindEquipementByName returns the gears of the armor set name.
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
