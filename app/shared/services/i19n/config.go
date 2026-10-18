package i19n

import (
	"fmt"
	plumtypes "github.com/cyril-fl/plum/types"
	"os"
	"slices"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// - Translator - //
type Translator struct {
	bundle       *i18n.Bundle
	local        *i18n.Localizer
	localcatalog *i18n.Localizer
}

var cache *Translator

func New() *Translator {
	return &Translator{
		bundle:       nil,
		local:        nil,
		localcatalog: nil,
	}
}

func (t *Translator) SetBundle(b *i18n.Bundle) {
	t.bundle = b
}

func (t *Translator) SetLocale(l *i18n.Localizer) {
	t.local = l
}

func (t *Translator) SetLocalerror(l *i18n.Localizer) {
	t.localcatalog = l
}

func SetCache(i19n *Translator) {
	cache = i19n
}

// - Locale - //
type Locale struct {
	Code string `yaml:"code"`
	ISO  string `yaml:"iso"`
	Name string `yaml:"name"`
	File string `yaml:"file"`
}

func (locale *Locale) ParseTag() (language.Tag, error) {
	tag, err := language.Parse(locale.ISO)
	if err != nil {
		return language.Tag{}, fmt.Errorf("invalid locale '%s': %w", locale.ISO, err)
	}

	return tag, nil
}

// - format - //
type format string

const (
	FORMAT_JSON format = "json"
	FORMAT_YAML format = "yaml"
)

var FormatValidExtensions = map[format][]string{
	FORMAT_JSON: {".json"},
	FORMAT_YAML: {".yaml", ".yml"},
}

// - Config - //
type Config struct {
	Default string   `yaml:"default"`
	Dir     string   `yaml:"dir"`
	Format  format   `yaml:"format"`
	Locales []Locale `yaml:"locales"`
}

func (cfg Config) GetDefaultLocale(defaultLocale string) (*Locale, error) {
	return cfg.GetLocaleByCode(defaultLocale)
}

func (cfg Config) GetLocaleByCode(code string) (*Locale, error) {
	for index := range cfg.Locales {
		if cfg.Locales[index].Code == code {
			return &cfg.Locales[index], nil
		}
	}

	return nil, fmt.Errorf("locale '%s' not found", code)
}

func (cfg Config) GetLocaleByISO(iso string) (*Locale, error) {
	for index := range cfg.Locales {
		if cfg.Locales[index].ISO == iso {
			return &cfg.Locales[index], nil
		}
	}

	return nil, fmt.Errorf("locale '%s' not found", iso)
}

func (cfg Config) Validate() error {
	if err := isDirectoryValid(cfg.Dir); err != nil {
		return err
	}

	if err := isFileFormatValid(cfg.Format); err != nil {
		return err
	}

	if err := isLocalesConfigValid(cfg.Default, cfg.Locales); err != nil {
		return err
	}

	return nil
}

// - Validation - //
func isDirectoryValid(path string) error {
	if path == "" {
		return fmt.Errorf("i18n directory is not specified")
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("check i18n directory '%s': %w", path, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("i18n directory '%s' does not exist", path)
	}

	return nil
}

func isFileFormatValid(extension format) error {
	if !slices.Contains([]format{FORMAT_JSON, FORMAT_YAML}, extension) {
		return fmt.Errorf("invalid i18n format '%s', must be 'json' or 'yaml'", extension)
	}

	return nil
}

func isLocalesConfigValid(defaultLocale string, locales []Locale) error {
	localeSet := make(plumtypes.Set[string])

	for _, locale := range locales {
		if localeSet.Has(locale.Code) {
			return fmt.Errorf("duplicate locale code '%s'", locale.Code)
		}

		if err := isLocaleValid(locale); err != nil {
			return err
		}

		localeSet.Add(locale.Code)
	}

	if !localeSet.Has(defaultLocale) {
		return fmt.Errorf("default locale '%s' is not in the list of available locales", defaultLocale)
	}

	return nil
}

func isLocaleValid(locale Locale) error {
	if locale.Code == "" {
		return fmt.Errorf("locale code is not specified for one of the locales")
	}
	if locale.ISO == "" {
		return fmt.Errorf("locale ISO code is not specified for locale '%s'", locale.Code)
	}
	if locale.Name == "" {
		return fmt.Errorf("locale name is not specified for locale '%s'", locale.Code)
	}
	if locale.File == "" {
		return fmt.Errorf("locale file is not specified for locale '%s'", locale.Code)
	}
	return nil
}
