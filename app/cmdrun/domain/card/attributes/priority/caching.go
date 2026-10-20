package priority

import (
	"solopg/app/shared/services/logs"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name   string `yaml:"name"`
	Weight int    `yaml:"weight"`
}

var filePath = config.Current.Documents.Files.Priorities

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
	for _, priority := range List() {
		if priority.GetName() == value.Name {
			return true
		}
	}
	return false
}

func AddInConfig(value Template) error {
	priorities, err := ListWithError()
	if err != nil {
		reportLoadError()
		return err
	}

	for _, priority := range priorities {
		if priority.GetName() == value.Name {
			return logs.MildError("error.invalid:already_exists", map[string]any{
				"Subject": "priority",
				"Value":   value.Name,
			})
		}
	}

	cachedConfig.Update(func(config *[]yamlConfig) {
		*config = append(*config, yamlConfig{
			Name:   value.Name,
			Weight: value.Weight,
		})
	})

	return nil
}

func ListWithError() ([]Priority, error) {
	configs, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil, err
	}

	return makeSet(configs), nil
}

func List() []Priority {
	priorities, err := ListWithError()
	if err != nil {
		reportLoadError()
		return nil
	}

	return priorities
}

// MakeDefault adds a priority to the runtime cache when it is not already
// configured and returns the resulting value.
func MakeDefault(defaultValue string) Priority {
	_, _ = ListWithError()

	if existing := FindByName(defaultValue); existing != nil {
		return *existing
	}

	cachedConfig.Update(func(config *[]yamlConfig) {
		*config = append(*config, yamlConfig{Name: defaultValue})
	})

	return *New(Template{Name: defaultValue})
}

func reportLoadError() {
	cachedConfig.ReportError(func(err error) {
		logs.Error("error.loading:file", map[string]any{
			"Subject": "priority",
			"File":    filePath,
			"Error":   err,
		})
	})
}

// - Collection - //
func ListNames() []string {
	priorities := List()
	list := make([]string, 0, len(priorities))
	for _, priority := range priorities {
		list = append(list, priority.GetName())
	}

	return list
}

func FindByName(name string) *Priority {
	for _, priority := range List() {
		if priority.GetName() == name {
			return &priority
		}
	}

	return nil
}

func makeSet(configs []yamlConfig) []Priority {
	set := make([]Priority, 0, len(configs))

	for _, config := range configs {
		set = append(set, newFromYaml(config))
	}

	return set
}

func newFromYaml(config yamlConfig) Priority {
	return Priority{
		name:   config.Name,
		weight: config.Weight,
	}
}
