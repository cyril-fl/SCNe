package hint

import (
	"errors"

	"solopg/app/cmdrun/domain/gameplay/dice"
	"solopg/app/shared/services/logs"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillentries "github.com/cyril-fl/plum/services/factory/millentries"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name   string   `yaml:"name"`
	Values []string `yaml:"values"`
}

var folderConfigPath = config.GetConfig().Documents.Folders.Hint

var cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

func loadFromSource() ([]yamlConfig, error) {
	entries := plummillentries.New()
	files := entries.GetFilesFromSource(folderConfigPath, true)

	if entries.HasErr() {
		return nil, logs.RawError("load hint folder %s: %w", folderConfigPath, entries.GetErr())
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
		return yamlConfig{}, logs.RawError("load hint file %s: %w", fileAddress, loader.GetErr())
	}

	return *params, nil
}

// - Hint - //

// Hint is a named list of hint words.
type Hint struct {
	name   string
	values []string
}

func newFromYamlConfig(config yamlConfig) Hint {
	return Hint{
		name:   config.Name,
		values: append([]string(nil), config.Values...),
	}
}

// GetByName returns the hint named name, or an error.
func GetByName(name string) (*Hint, error) {
	hints := List()

	for _, hint := range hints {
		if hint.name == name {
			return &hint, nil
		}
	}

	return nil, logs.Error("error.not_found.id", map[string]any{
		"Subject": "hint",
		"ID":      name,
	})
}

// GetName returns the hint name.
func (o Hint) GetName() string {
	return o.name
}

// Values returns a copy of the hint words.
func (o Hint) Values() []string {
	return append([]string(nil), o.values...)
}

// - Collection - //

// ListWithError returns the hints loaded from the hint folder, once.
func ListWithError() ([]Hint, error) {
	configs, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return nil, err
	}
	return makeSet(configs), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
func List() []Hint {
	hints, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.SilentWarning("error.loading", map[string]any{
				"Subject": "hint",
				"Error":   err,
			})
		})
		return nil
	}
	return hints
}

func makeSet(configs []yamlConfig) []Hint {
	result := make([]Hint, 0, len(configs))

	for _, config := range configs {
		result = append(result, newFromYamlConfig(config))
	}

	return result
}

// - Helpers - //

// Roll returns three random words of hint, repeats allowed, or nil when it has none.
func Roll(hint Hint) []string {
	res := []string{}
	values := hint.Values()
	if len(values) == 0 {
		return nil
	}

	for i := 0; i < 3; i++ {
		value := dice.Roll(len(values))
		res = append(res, values[value-1])
	}

	return res
}
