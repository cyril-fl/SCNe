package config

import (
	"fmt"
	"os"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"strings"

	"charm.land/lipgloss/v2"
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
	assertConfigFile(file)

	loader := plummillload.New(file)
	loaded := loader.LoadFromFile[config]()

	if loader.HasErr() {
		panic(logs.RawError("%w", loader.GetErr()))
	}

	if err := loaded.I18n.Validate(); err != nil {
		panic(logs.Error("error.invalid", map[string]any{
			"Subject": "config",
			"Error":   err,
		}))
	}

	current = *loaded
}

// assertConfigFile exits when CONFIG_FILE is unset or missing.
// It runs before logs.Init and the DB connection, so the error is only printed.
func assertConfigFile(file string) {
	if file == "" {
		exitWithLog("CONFIG_FILE is not set. Create a .env file (see .env.example) or export CONFIG_FILE.")
	}

	if _, err := os.Stat(file); err != nil {
		exitWithLog(fmt.Sprintf("config file %q not found (CONFIG_FILE). Check your .env file.", file))
	}
}

func exitWithLog(message string) {
	lipgloss.Fprintln(os.Stderr, logs.New(logs.Template{Type: logs.ERR, Message: message}).String())
	os.Exit(1)
}

func init() {
	godotenv.Load()
	load()
}
