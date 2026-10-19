package config

import (
	"os"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"strings"

	"github.com/joho/godotenv"
)

type config struct {
	Name      string            `yaml:"name"`
	Verbose   bool              `yaml:"verbose"`
	I18n      i19n.Config       `yaml:"i18n"`
	Commands  commands          `yaml:"commands"`
	Documents structureDocument `yaml:"documents"`
}

var Current config

func load() (*config, error) {
	file := strings.TrimSpace(os.Getenv("CONFIG_FILE"))
	current, err := yaml.LoadFromFile[config](file)

	if err != nil {
		return nil, err
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
