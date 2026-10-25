package config

import (
	_ "embed"
	"fmt"
	"os"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/theme"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// privateFile is the built-in config, used as fallback for every key CONFIG_FILE leaves out.
//
//go:embed .config.yaml
var privateFile []byte

type config struct {
	Name     string   `yaml:"name"`
	Commands commands `yaml:"commands"`
	Public   `yaml:",inline"`
}

// Public holds what CONFIG_FILE may override.
type Public struct {
	I18n      i19n.Config       `yaml:"i18n"`
	Documents structureDocument `yaml:"documents"`
	Theme     theme.Config      `yaml:"theme"`
}

var current config

func GetConfig() config {
	return current
}

func load() {
	var loaded config
	if err := yaml.Unmarshal(privateFile, &loaded); err != nil {
		panic(logs.RawError("%w", err))
	}

	if file := strings.TrimSpace(os.Getenv("CONFIG_FILE")); file != "" {
		loadPublic(&loaded, file)
	}

	if err := loaded.I18n.Validate(); err != nil {
		panic(logs.Error("error.invalid", map[string]any{
			"Subject": "config",
			"Error":   err,
		}))
	}

	current = loaded
}

// loadPublic decodes file over the public part of cfg; the locales are appended to the private ones.
func loadPublic(cfg *config, file string) {
	assertConfigFile(file)

	data, err := os.ReadFile(file)
	if err != nil {
		panic(logs.RawError("%w", err))
	}

	privateLocales := slices.Clone(cfg.I18n.Locales)
	if err := yaml.Unmarshal(data, &cfg.Public); err != nil {
		panic(logs.RawError("%w", err))
	}

	// TODO permettre au joueur d'ajouter ses propres fichiers de langue (aujourd'hui ils doivent être dans i18n.dir)
	cfg.I18n.Locales = mergeLocales(privateLocales, cfg.I18n.Locales)
}

// mergeLocales appends public to private; a public locale replaces the private one with the same code.
func mergeLocales(private []i19n.Locale, public []i19n.Locale) []i19n.Locale {
	merged := slices.Clone(private)

	for _, locale := range public {
		index := slices.IndexFunc(merged, func(l i19n.Locale) bool { return l.Code == locale.Code })
		if index < 0 {
			merged = append(merged, locale)
			continue
		}
		merged[index] = locale
	}

	return merged
}

// assertConfigFile exits when CONFIG_FILE is missing.
// It runs before logs.Init and the DB connection, so the error is only printed.
func assertConfigFile(file string) {
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
