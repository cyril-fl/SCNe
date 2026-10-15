package hint

import (
	"errors"
	"path/filepath"

	"solopg/app/cmdrun/domain/gameplay/dice"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"solopg/config"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name   string   `yaml:"name"`
	Values []string `yaml:"values"`
}

var folderConfigPath = config.Current.Documents.Folders.Hint

var cachedConfig []yamlConfig

func loadFromSource() error {
	//  LOW ameliorer. Si c'est un fichier on load si c'est un folder on load mais autrement
	// Appliquer cette logique la au oracles dice ect...
	files, err := yaml.GetFilesFromSource(folderConfigPath, true)
	if err != nil {
		return logs.Error("error.loading:folder", map[string]any{
			"Subject": "hint",
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
				"Subject": "hint",
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
		return logs.Error("error.loading:file", map[string]any{
			"Subject": "hint",
			"File":    fileAddress,
			"Error":   err,
		})
	}

	cachedConfig = append(cachedConfig, *params)

	return nil
}

// - Hint - //
type Hint struct {
	name   string
	values []string
}

func newFromYamlConfig(config yamlConfig) Hint {
	return Hint{
		name:   config.Name,
		values: config.Values,
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
	return o.values
}

// - Collection - //
func List() []Hint {
	if len(cachedConfig) == 0 {
		err := loadFromSource()
		if err != nil {
			return nil
		}
	}

	return makeSet(cachedConfig)
}

func ListFromFolder(path string) ([]string, error) {
	list, err := yaml.GetFolderFiles(path)
	if err != nil {

		return nil, logs.Error("error.loading:folder", map[string]any{
			"Subject": "hint",
			"Folder":  path,
			"Error":   err,
		})
	}

	errs := []error{}
	hints := []string{}
	for _, file := range list {
		if params, err := yaml.LoadFromFile[yamlConfig](filepath.Join(path, file)); err != nil {

			errs = append(errs, logs.Error("error.loading:file", map[string]any{
				"Subject": "hint",
				"File":    file,
				"Error":   err,
			}))
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
