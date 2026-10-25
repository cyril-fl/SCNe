package oracle

import (
	"errors"
	"fmt"
	"path/filepath"

	"scne/app/cmdrun/domain/gameplay/dice"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"
	"scne/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillentries "github.com/cyril-fl/plum/services/factory/millentries"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"

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

// Result is an oracle roll and the result of its interval.
type Result[T any] struct {
	Roll     int  `yaml:"roll"`
	Critical bool `yaml:"critical"`
	Result   T    `yaml:"result"`
}

var folderConfigPath = config.GetConfig().Documents.Folders.Oracle

var cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

func loadFromSource() ([]yamlConfig, error) {
	entries := plummillentries.New()
	files := entries.GetFilesFromSource(folderConfigPath, true)
	if entries.HasErr() {
		return nil, logs.RawError("load oracle folder %s: %w", folderConfigPath, entries.GetErr())
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
		return yamlConfig{}, logs.RawError("load oracle file %s: %w", fileAddress, loader.GetErr())
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

// Oracle maps dice roll intervals to results.
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

// GetByName returns the oracle named name, or an error.
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

// GetName returns the oracle translation key, "oracles.<name>".
func (o Oracle) GetName() string {
	return fmt.Sprintf("oracles.%s", o.name)
}

// IsVisible reports whether the oracle is shown to the player.
func (o Oracle) IsVisible() bool {
	return o.visible
}

// - Collection - //

// ListWithError returns the oracles loaded from the oracle folder, once.
func ListWithError() ([]Oracle, error) {
	configs, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return nil, err
	}
	return makeSet(configs), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
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

// ListVisible returns the visible oracles.
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

// ListFromFolder returns the oracle names of the files in path, uncached.
func ListFromFolder(path string) ([]string, error) {
	entries := plummillentries.New()
	list := entries.GetFolderFiles(path)
	if entries.HasErr() {

		return nil, logs.RawError("load oracle folder %s: %w", path, entries.GetErr())
	}

	errs := []error{}
	oracles := []string{}
	for _, file := range list {
		loader := plummillload.New(filepath.Join(path, file))
		params := loader.LoadFromFile[yamlConfig]()
		if loader.HasErr() {

			errs = append(errs, logs.RawError("load oracle file %s: %w", file, loader.GetErr()))
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

// Roll rolls o and returns the matching interval result as T.
// It fails when no interval matches or the result is not a T.
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
