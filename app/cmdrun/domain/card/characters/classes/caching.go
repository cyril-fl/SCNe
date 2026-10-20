package classes

import (
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/shared/services/logs"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name     string           `yaml:"name"`
	Playable bool             `yaml:"playable"`
	Bonus    []stats.Modifier `yaml:"bonus"`
	ArmorSet string           `yaml:"armor_set"`
}

var filePath = config.Current.Documents.Files.Classes

var cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

func loadFromFile() ([]yamlConfig, error) {
	loader := plummillload.New(filePath)
	params := loader.LoadListFromFile[yamlConfig]()
	if loader.HasErr() {
		return nil, logs.RawError("%w", loader.GetErr())
	}

	return params, nil
}

func Has(value Template) bool {
	new := New(value)

	for _, class := range List() {
		if class.GetHash() == new.GetHash() {
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
			"Subject": "class",
			"Value":   value.Name,
		})
	}

	new := yamlConfig{
		Name:     value.Name,
		Playable: value.Playable,
		Bonus:    value.Bonus,
		ArmorSet: value.EquipementName,
	}

	cachedConfig.Update(func(config *[]yamlConfig) {
		*config = append(*config, new)
	})

	return nil
}

// - Collection - //
func ListWithError() ([]Class, error) {
	configs, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}

	return makeSet(configs), nil
}

func List() []Class {
	classes, err := ListWithError()
	if err != nil {
		reportLoadError()
		return nil
	}

	return classes
}

func reportLoadError() {
	cachedConfig.ReportError(func(err error) {
		logs.SilentWarning("error.loading", map[string]any{
			"Subject": "class",
			"File":    filePath,
			"Error":   err,
		})
	})
}

func ListNames() []string {
	var list []string
	for _, i := range List() {
		list = append(list, i.GetName())
	}
	return list
}

func FindByName(name string) *Class {
	for _, i := range List() {
		if i.GetName() == name {
			return &i
		}
	}
	return nil
}

func InsertIfIsnt(name string) (*Class, error) {
	if class := FindByName(name); class != nil {
		return class, nil
	}

	if err := AddInConfig(Template{
		Name:     name,
		Playable: true,
	}); err != nil {
		return nil, err
	}

	return FindByName(name), nil
}

func makeSet(configs []yamlConfig) []Class {
	result := make([]Class, 0, len(configs))

	for _, config := range configs {
		result = append(result, newFromYaml(config))
	}

	return result
}

func newFromYaml(config yamlConfig) Class {
	return Class{
		name:           config.Name,
		playable:       config.Playable,
		bonus:          config.Bonus,
		equipementName: config.ArmorSet,
	}
}
