package combat

import (
	"errors"
	"slices"

	"solopg/app/cmdrun/domain/gameplay/oracle"
	"solopg/app/shared/services/logs"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillentries "github.com/cyril-fl/plum/services/factory/millentries"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	ID      string                    `yaml:"id"`
	Oracle  string                    `yaml:"oracle"`
	Effects map[string]map[string]int `yaml:"effects"`
}

const (
	effectDamage = "damage"
	effectGuard  = "guard"
)

var supportedEffects = []string{effectDamage, effectGuard}

var folderConfigPath = config.GetConfig().Documents.Folders.Actions

var cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

func loadFromSource() ([]yamlConfig, error) {
	entries := plummillentries.New()
	files := entries.GetFilesFromSource(folderConfigPath, true)
	if entries.HasErr() {
		return nil, logs.RawError("load action folder %s: %w", folderConfigPath, entries.GetErr())
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
		return yamlConfig{}, logs.RawError("load action file %s: %w", fileAddress, loader.GetErr())
	}

	if err := validate(*params); err != nil {
		return yamlConfig{}, logs.RawError("load action file %s: %w", fileAddress, err)
	}

	return *params, nil
}

func validate(config yamlConfig) error {
	if config.ID == "" {
		return logs.RawError("action has no id")
	}

	if _, err := oracle.GetByName(config.Oracle); err != nil {
		return logs.RawError("action %s: %w", config.ID, err)
	}

	for outcome, effect := range config.Effects {
		for kind, value := range effect {
			if !slices.Contains(supportedEffects, kind) {
				return logs.RawError("action %s, outcome %s: unknown effect %q (supported: %v)", config.ID, outcome, kind, supportedEffects)
			}
			if value < 0 {
				return logs.RawError("action %s, outcome %s: effect %s must be >= 0, got %d", config.ID, outcome, kind, value)
			}
		}
	}

	return nil
}

// - Action - //

// Action is a combat action resolved by its oracle; the oracle result selects the effect.
type Action struct {
	id      string
	oracle  string
	effects map[string]map[string]int
}

func newFromYamlConfig(config yamlConfig) Action {
	return Action{
		id:      config.ID,
		oracle:  config.Oracle,
		effects: config.Effects,
	}
}

// GetID returns the action ID.
func (a Action) GetID() string {
	return a.id
}

// GetName returns the action translation key, "actions.<id>".
func (a Action) GetName() string {
	return "actions." + a.id
}

// GetByID returns the action with id, or an error.
func GetByID(id string) (*Action, error) {
	actions, err := ListWithError()
	if err != nil {
		return nil, err
	}

	for _, action := range actions {
		if action.id == id {
			return &action, nil
		}
	}

	return nil, logs.Error("error.not_found.id", map[string]any{
		"Subject": "action",
		"ID":      id,
	})
}

// - Collection - //

// ListWithError returns the actions loaded from the action folder, once.
func ListWithError() ([]Action, error) {
	configs, err := cachedConfig.Load(loadFromSource)
	if err != nil {
		return nil, err
	}
	return makeSet(configs), nil
}

// List is ListWithError that logs the load error once and returns nil on failure.
func List() []Action {
	actions, err := ListWithError()
	if err != nil {
		cachedConfig.ReportError(func(err error) {
			logs.SilentError("error.loading", map[string]any{
				"Subject": "action",
				"Error":   err,
			})
		})
		return nil
	}
	return actions
}

func makeSet(configs []yamlConfig) []Action {
	result := make([]Action, 0, len(configs))

	for _, config := range configs {
		result = append(result, newFromYamlConfig(config))
	}

	return result
}
