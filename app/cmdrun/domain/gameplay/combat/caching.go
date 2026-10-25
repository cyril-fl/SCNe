package combat

import (
	"solopg/app/shared/services/logs"
	"solopg/config"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

// - Configuration & caching - //
type yamlConfig struct {
	Shortcuts []string `yaml:"shortcuts"`
}

var fileConfigPath = config.GetConfig().Documents.Files.Combat

var cachedConfig = plumprimitive.NewCache[yamlConfig]()

func loadFromFile() (yamlConfig, error) {
	loader := plummillload.New(fileConfigPath)
	params := loader.LoadFromFile[yamlConfig]()
	if loader.HasErr() {
		return yamlConfig{}, logs.RawError("load combat file %s: %w", fileConfigPath, loader.GetErr())
	}

	return *params, nil
}

// - Shortcuts - //

// GetShortcuts returns the keys rolling the Nth stat besides alt+N, or nil when the file fails to load.
func GetShortcuts() []string {
	config, err := cachedConfig.Load(loadFromFile)
	if err != nil {
		return nil
	}
	return config.Shortcuts
}
