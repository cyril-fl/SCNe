package config

import (
	"os"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"strings"

	"github.com/joho/godotenv"

	plummillload "github.com/cyril-fl/plum/services/factory/millload"
)

type config struct {
	Name      string            `yaml:"name"`
	I18n      i19n.Config       `yaml:"i18n"`
	Commands  commands          `yaml:"commands"`
	Documents structureDocument `yaml:"documents"`
}

var Current config

func load() (*config, error) {
	file := strings.TrimSpace(os.Getenv("CONFIG_FILE"))

	loader := plummillload.New(file)
	current := loader.LoadFromFile[config]()
	if loader.HasErr() {
		return nil, logs.RawError("%w", loader.GetErr())
	}

	if err := current.I18n.Validate(); err != nil {
		return nil, logs.Error("error.invalid", map[string]any{
			"Subject": "config",
			"Error":   err,
		})
	}

	return current, nil
}

func init() {
	godotenv.Load()
	loaded, err := load()
	if err != nil {
		panic(err)
	}
	Current = *loaded
}
