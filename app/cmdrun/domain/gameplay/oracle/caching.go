package oracle

import (
	"errors"
	"fmt"
	"path/filepath"

	"solopg/app/cmdrun/domain/gameplay/dice"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// - Configuration & caching - //
type yamlConfig struct {
	Name      string     `yaml:"name"`
	Dice      int        `yaml:"dice"`
	Type      string     `yaml:"type"`
	Visible   bool       `yaml:"visible"`
	Intervals []interval `yaml:"intervals"`
}

type interval struct {
	Min      int  `yaml:"min"`
	Max      int  `yaml:"max"`
	Critical bool `yaml:"critical"`
	Result   any  `yaml:"result"`
}

type Result[T any] struct {
	Roll     int  `yaml:"roll"`
	Critical bool `yaml:"critical"`
	Result   T    `yaml:"result"`
}

var folderConfigPath = config.Current.Documents.Folders.Oracle

var cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

func loadFromSource() ([]yamlConfig, error) {
	files, err := yaml.GetFilesFromSource(folderConfigPath, true)
	if err != nil {
		return nil, fmt.Errorf("load oracle folder %s: %w", folderConfigPath, err)
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
		return yamlConfig{}, fmt.Errorf("load oracle file %s: %w", fileAddress, err)
	}

	if params.Dice <= 0 {
		return yamlConfig{}, fmt.Errorf("oracle %s has an invalid dice count", params.Name)
	}

	if len(params.Intervals) == 0 {
		return yamlConfig{}, fmt.Errorf("oracle %s has no intervals", params.Name)
	}

	return *params, nil
}

// - Oracle - //
type Oracle struct {
	name      string
	dice      int
	typ       string
	visible   bool
	intervals []interval
}

func newFromYamlConfig(config yamlConfig) Oracle {
	return Oracle{
		name:      config.Name,
		dice:      config.Dice,
		typ:       config.Type,
		visible:   config.Visible,
		intervals: append([]interval(nil), config.Intervals...),
	}
}

func GetByName(name string) (*Oracle, error) {
	oracles := List()

	for _, oracle := range oracles {
		if oracle.name == name {
			return &oracle, nil
		}
	}

	return nil, logs.Error("error.not_found.id", map[string]any{
		"Subject": plumtransform.Capitalize(i19n.Localize("oracle")),
		"ID":      name,
	})
}

func (o Oracle) GetName() string {
	return fmt.Sprintf("oracles.%s", o.name)
}

func (o Oracle) IsVisible() bool {
	return o.visible
}

// - Collection - //
func ListWithError() ([]Oracle, error) {
	configs, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return nil, err
	}
	return makeSet(configs), nil
}

func List() []Oracle {
	oracles, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.Error("error.loading", map[string]any{
				"Subject": plumtransform.Capitalize(i19n.Localize("oracle")),
				"Error":   err,
			})
		})
		return nil
	}
	return oracles
}

func ListVisible() []Oracle {
	oracles := List()
	visibleOracles := make([]Oracle, 0)

	for _, oracle := range oracles {
		if oracle.IsVisible() {
			visibleOracles = append(visibleOracles, oracle)
		}
	}
	return visibleOracles
}

func ListFromFolder(path string) ([]string, error) {
	list, err := yaml.GetFolderFiles(path)
	if err != nil {

		return nil, logs.Error("error.loading:folder", map[string]any{
			"Subject": plumtransform.Capitalize(i19n.Localize("oracle")),
			"Folder":  path,
			"Error":   err,
		})
	}

	errs := []error{}
	oracles := []string{}
	for _, file := range list {
		if params, err := yaml.LoadFromFile[yamlConfig](filepath.Join(path, file)); err != nil {

			errs = append(errs, logs.Error("error.loading:file", map[string]any{
				"Subject": plumtransform.Capitalize(i19n.Localize("oracle")),
				"File":    file,
				"Error":   err,
			}))
		} else {
			oracles = append(oracles, params.Name)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return oracles, nil
}

func makeSet(configs []yamlConfig) []Oracle {
	result := make([]Oracle, 0, len(configs))

	for _, config := range configs {
		result = append(result, newFromYamlConfig(config))
	}

	return result
}

// - Helpers - //
func Roll[T any](o Oracle) (*Result[T], error) {
	value := dice.Roll(o.dice)
	for _, interval := range o.intervals {
		if value >= interval.Min && value <= interval.Max {
			res, ok := interval.Result.(T)
			if !ok {

				return nil, logs.Error("error.unexpected:value", map[string]any{
					"Subject":  o.name,
					"Expected": fmt.Sprintf("%T", res),
					"Value":    fmt.Sprintf("%T", interval.Result),
				})
			}

			return &Result[T]{
				Roll:     value,
				Critical: interval.Critical,
				Result:   res,
			}, nil
		}
	}

	return nil, logs.Error("error.not_found.id", map[string]any{
		"Subject": o.name,
		"ID":      value,
	})
}
