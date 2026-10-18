package initI19n

import (
	"fmt"
	"path/filepath"
	"slices"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"strings"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	y "gopkg.in/yaml.v3"
)

type run struct {
	plumprimitive.Process
	cache

	cfg    i19n.Config
	locale string
}

type cache struct {
	defaultLocale *i19n.Locale

	bundle     *i18n.Bundle
	local      *i18n.Localizer
	localerror *i18n.Localizer

	files []string
}

func Process(cfg i19n.Config, locale string) *run {
	return &run{
		cfg:    cfg,
		locale: locale,
	}
}

func (p *run) Run() {
	// Set env
	p.reset()
	p.assertConfig()
	p.setDefaultLocale()

	// Bundle
	p.initBundle()
	p.loadLocaleFiles()
	p.registerLocaleFiles()

	// Localizer
	p.initLocalizer()
	p.initErrorLocalizer()

	// Set
	p.seti19n()
}

// Getters & Setters
func (p *run) GetResult() {
	logs.SilentWarning("error:not_implemented", map[string]any{
		"Function": "GetResult",
		"Subject":  "init:i19n",
	})
}

// Methods
func (p *run) assertConfig() {
	if p.HasErr() {
		return
	}

	if err := p.cfg.Validate(); err != nil {
		p.SetErr(err)
	}
}

func (p *run) setDefaultLocale() {
	if p.HasErr() {
		return
	}

	defaultLocale, err := p.cfg.GetDefaultLocale(p.cfg.Default)
	if err != nil {
		p.SetErr(err)
		return
	}

	p.defaultLocale = defaultLocale
}

func (p *run) initBundle() {
	if p.HasErr() {
		return
	}

	tag, err := p.defaultLocale.ParseTag()
	if err != nil {
		p.SetErr(err)
		return
	}

	bundle := i18n.NewBundle(tag)
	bundle.RegisterUnmarshalFunc(string(p.cfg.Format), y.Unmarshal)

	p.bundle = bundle
}

func (p *run) loadLocaleFiles() {
	if p.HasErr() {
		return
	}

	files, err := loadLocaleFile(p.cfg)
	if err != nil {

		p.SetErr(err)
		return
	}

	p.files = files
}

func (p *run) registerLocaleFiles() {
	if p.HasErr() {
		return
	}

	for _, path := range p.files {
		if _, err := p.bundle.LoadMessageFile(path); err != nil {
			p.SetErr(fmt.Errorf("load message file '%s': %w", path, err))
		}
	}
}

func (p *run) initLocalizer() {
	if p.HasErr() {
		return
	}

	localizer, err := newLocalizer(p.cfg, p.bundle, p.locale, p.defaultLocale)
	if err != nil {
		p.SetErr(err)
		return
	}

	p.local = localizer
}

func (p *run) initErrorLocalizer() {
	if p.HasErr() {
		return
	}

	localizer, err := newLocalizer(p.cfg, p.bundle, "en", p.defaultLocale)
	if err != nil {
		p.SetErr(err)
		return
	}

	p.localerror = localizer
}

func (p *run) seti19n() {
	if p.HasErr() {
		return
	}

	t := i19n.New()
	t.SetBundle(p.bundle)
	t.SetLocale(p.local)
	t.SetLocalerror(p.localerror)

	i19n.SetCache(t)
}

func (p *run) reset() {
	p.cache = cache{}

	p.SetErr(nil)
	i19n.SetCache(i19n.New())
}

// Helpers
func loadLocaleFile(cfg i19n.Config) ([]string, error) {
	files, err := yaml.GetFilesFromSource(cfg.Dir, true)
	if err != nil {
		return nil, err
	}

	assertedFiles := make([]string, 0)
	for _, path := range files {
		extensions := i19n.FormatValidExtensions[cfg.Format]

		if !(slices.Contains(extensions, filepath.Ext(path))) {
			continue
		}

		name := filepath.Base(path)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		if _, err := cfg.GetLocaleByISO(name); err != nil {
			continue
		}

		assertedFiles = append(assertedFiles, path)
	}

	return assertedFiles, nil
}

/*
NOTE "newLocalizer" may very well be a method, but I find that it mixes responsibilities.
Thus, a localizer can be initialized without going through the initializer.
*/
func newLocalizer(cfg i19n.Config, bundle *i18n.Bundle, lang string, defaultLocale *i19n.Locale) (*i18n.Localizer, error) {
	locale, err := getLocale(cfg, lang, defaultLocale)
	if err != nil {
		return nil, err
	}

	tag, err := locale.ParseTag()
	if err != nil {
		return nil, err
	}

	defaultTag, err := defaultLocale.ParseTag()
	if err != nil {
		return nil, err
	}

	return i18n.NewLocalizer(bundle, tag.String(), defaultTag.String()), nil
}

func getLocale(cfg i19n.Config, lang string, defaultLocale *i19n.Locale) (*i19n.Locale, error) {
	if lang == "" {
		return defaultLocale, nil
	}

	locale, err := cfg.GetLocaleByCode(lang)
	if err != nil {
		return nil, fmt.Errorf("locale '%s' not found: %w", lang, err)
	}

	return locale, nil
}
