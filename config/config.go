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

var current config

func GetConfig() config {
	return current
}

func load() {
	file := strings.TrimSpace(os.Getenv("CONFIG_FILE"))

	loader := plummillload.New(file)
	loaded := loader.LoadFromFile[config]()

	var e error

	if loader.HasErr() {
		e = logs.RawError("%w", loader.GetErr())
	}

	if err := loaded.I18n.Validate(); err != nil {
		e = logs.Error("error.invalid", map[string]any{
			"Subject": "config",
			"Error":   err,
		})
		panic(e)
	}

	if e != nil {
		panic(e)
	}

	current = *loaded
}

func init() {
	godotenv.Load()
	load()
}
