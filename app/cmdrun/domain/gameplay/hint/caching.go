package hint

import (
	"errors"
	"path/filepath"

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

var folderConfigPath = config.Current.Documents.Folders.Hint

var cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

func loadFromSource() ([]yamlConfig, error) {
	//  LOW ameliorer. Si c'est un fichier on load si c'est un folder on load mais autrement
	// Appliquer cette logique la au oracles dice ect...
	entries := plummillentries.New()
	files := entries.GetFilesFromSource(folderConfigPath, true)
	if entries.HasErr() {
		return nil, logs.RawError("load hint folder %s: %w", folderConfigPath, entries.GetErr())
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
	loader := plummillload.New(fileAddress)
	params := loader.LoadFromFile[yamlConfig]()
	if loader.HasErr() {
		return yamlConfig{}, logs.RawError("load hint file %s: %w", fileAddress, loader.GetErr())
	}
	return *params, nil
}

// - Hint - //
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

func (o Hint) GetName() string {
	return o.name
}
func (o Hint) Values() []string {
	return append([]string(nil), o.values...)
}

// - Collection - //
func ListWithError() ([]Hint, error) {
	configs, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return nil, err
	}
	return makeSet(configs), nil
}

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

func ListFromFolder(path string) ([]string, error) {
	entries := plummillentries.New()
	list := entries.GetFolderFiles(path)
	if entries.HasErr() {

		return nil, logs.RawError("load hint folder %s: %w", path, entries.GetErr())
	}

	errs := []error{}
	hints := []string{}
	for _, file := range list {
		loader := plummillload.New(filepath.Join(path, file))
		params := loader.LoadFromFile[yamlConfig]()
		if loader.HasErr() {

			errs = append(errs, logs.RawError("load hint file %s: %w", file, loader.GetErr()))
		} else {
			hints = append(hints, params.Name)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return hints, nil
}

func makeSet(configs []yamlConfig) []Hint {
	result := make([]Hint, 0, len(configs))

	for _, config := range configs {
		result = append(result, newFromYamlConfig(config))
	}

	return result
}

// - Helpers - //
func Roll(hint Hint) []string {
	res := []string{}
	values := hint.Values()

	for i := 0; i < 3; i++ {
		value := dice.Roll(len(values))
		res = append(res, values[value-1])
	}

	return res
}
