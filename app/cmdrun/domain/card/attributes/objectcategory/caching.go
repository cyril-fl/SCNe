package objectcategory

import (
	"slices"
	"solopg/app/shared/services/yaml"
	"solopg/config"
)

// - Configuration & caching - //
type yamlConfig struct {
	Values []Category `yaml:"values"`
}

var fileConfigPath = config.Current.Documents.Files.ObjectsCategory

var cachedConfig yamlConfig

func loadFromFile() error {
	params, err := yaml.LoadFromFile[yamlConfig](fileConfigPath)
	if err != nil {
		return err
	}

	cachedConfig = *params

	return nil
}

func (c *yamlConfig) addValue(value Category) {
	if !slices.Contains(c.Values, value) {
		c.Values = append(c.Values, value)
	}
}

// - Category - //
type Category string

func List() []Category {
	if len(cachedConfig.Values) == 0 {
		err := loadFromFile()
		if err != nil {
			return nil
		}
	}

	return cachedConfig.Values
}

/*
LOW Verrifier ou "MakeDefault" pourrait être utile
ex: Si la data vien d'une carte loader, c'est valider
*/
func MakeDefault(defaultValue string) Category {
	List()
	cachedConfig.addValue(Category(defaultValue))
	return Category(defaultValue)
}

func AssertWithDefault(provided string) Category {
	if Category(provided).Validate() {
		return Category(provided)
	}

	return MakeDefault("object")
}

func (r Category) Validate() bool {
	return slices.Contains(List(), r)
}

func (r Category) String() string {
	return string(r)
}
