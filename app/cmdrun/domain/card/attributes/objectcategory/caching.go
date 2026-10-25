package objectcategory

import (
	"slices"
	"solopg/app/shared/services/logs"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Values []Category `yaml:"values"`
}

var fileConfigPath = config.GetConfig().Documents.Files.ObjectsCategory

var cachedConfig = plumprimitive.NewCache[yamlConfig]()

func loadFromFile() (yamlConfig, error) {
	loader := plummillload.New(fileConfigPath)
	params := loader.LoadFromFile[yamlConfig]()
	if loader.HasErr() {
		return yamlConfig{}, logs.RawError("%w", loader.GetErr())
	}

	return *params, nil
}

func (c *yamlConfig) addValue(value Category) {
	if !slices.Contains(c.Values, value) {
		c.Values = append(c.Values, value)
	}
}

// - Category - //

// Category is a value from the objectcategory config file.
type Category string

// ListWithError returns a copy of the configured values; the file is loaded once.
func ListWithError() ([]Category, error) {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}

	return append([]Category(nil), config.Values...), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
func List() []Category {
	values, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.Error("error.loading:file", map[string]any{
				"Subject": "category",
				"File":    fileConfigPath,
				"Error":   err,
			})
		})
		return nil
	}

	return values
}

// MakeDefault adds defaultValue to the cached values when missing and returns it.
func MakeDefault(defaultValue string) Category {
	_, _ = ListWithError()
	cachedConfig.Update(func(config *yamlConfig) {
		config.addValue(Category(defaultValue))
	})
	return Category(defaultValue)
}

// AssertWithDefault returns provided when valid, otherwise the "object" category.
func AssertWithDefault(provided string) Category {
	if Category(provided).Validate() {
		return Category(provided)
	}

	return MakeDefault("object")
}

// Validate reports whether r is a configured value.
func (r Category) Validate() bool {
	return slices.Contains(List(), r)
}

// String returns r as a string.
func (r Category) String() string {
	return string(r)
}
