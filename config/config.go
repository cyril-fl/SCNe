package config

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/theme"
	"slices"

	"charm.land/lipgloss/v2"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// privateFile is the built-in config, used as fallback for every key the public file leaves out.
//
//go:embed .config.yaml
var privateFile []byte

// publicTemplate is written when the public file is created by `config edit`.
//
//go:embed .public.yaml
var publicTemplate []byte

type config struct {
	Name     string   `yaml:"name" json:"name"`
	Commands commands `yaml:"commands" json:"commands"`
	Public   `yaml:",inline"`
}

// Public holds what the public file may override.
type Public struct {
	I18n      i19n.Config       `yaml:"i18n"`
	Documents structureDocument `yaml:"documents" json:"documents"`
	Theme     theme.Config      `yaml:"theme" json:"theme"`
}

var current config

func GetConfig() config {
	return current
}

// PublicFile returns the path of the player's config file.
func PublicFile() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "scne", "config.yaml"), nil
}

// EnsurePublicFile creates path from the template when it does not exist.
func EnsurePublicFile(path string) error {
	if _, err := os.Stat(path); err == nil || !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, publicTemplate, 0o644)
}

// Check reports whether the public file loads over the private config and validates.
func Check() error {
	private, err := loadPrivate()
	if err != nil {
		return err
	}

	_, err = withPublic(private)
	return err
}

func load() {
	private, err := loadPrivate()
	if err != nil {
		panic(logs.RawError("%w", err))
	}

	if err := private.I18n.Validate(); err != nil {
		panic(logs.Error("error.invalid", map[string]any{
			"Subject": "config",
			"Error":   err,
		}))
	}

	current = private

	// An invalid public file is ignored so every command, `config edit` included, still runs.
	public, err := withPublic(private)
	if err != nil {
		printWarning(fmt.Sprintf("player config ignored: %v. Fix it with `scne config edit`.", err))
		return
	}

	current = public
}

func loadPrivate() (config, error) {
	var loaded config
	if err := yaml.Unmarshal(privateFile, &loaded); err != nil {
		return config{}, fmt.Errorf("unmarshal private config: %w", err)
	}

	return loaded, nil
}

// withPublic returns cfg overridden by the public file when it exists, then validated.
func withPublic(cfg config) (config, error) {
	// No resolvable config dir means no public file: the private config is used alone.
	file, err := PublicFile()
	if err != nil {
		return cfg, nil
	}

	if err := loadPublic(&cfg, file); err != nil {
		return config{}, err
	}

	if err := cfg.I18n.Validate(); err != nil {
		return config{}, fmt.Errorf("%s: %w", file, err)
	}

	return cfg, nil
}

// loadPublic decodes file over the public part of cfg; the locales are appended to the private ones.
// A missing file leaves cfg unchanged.
func loadPublic(cfg *config, file string) error {
	data, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	privateLocales := slices.Clone(cfg.I18n.Locales)
	if err := yaml.Unmarshal(data, &cfg.Public); err != nil {
		return fmt.Errorf("unmarshal %s: %w", file, err)
	}

	// TODO permettre au joueur d'ajouter ses propres fichiers de langue (aujourd'hui ils doivent être dans i18n.dir)
	cfg.I18n.Locales = mergeLocales(privateLocales, cfg.I18n.Locales)

	return nil
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

// printWarning prints message on stderr; it runs before logs.Init, so nothing is stored.
func printWarning(message string) {
	lipgloss.Fprintln(os.Stderr, logs.New(logs.Template{Type: logs.WARN, Message: message}).String())
}

func init() {
	godotenv.Load()
	load()
}
