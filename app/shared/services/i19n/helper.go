package i19n

import (
	"fmt"
	"strconv"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Localize translates id in the user locale; string data values are translated too.
// It returns "{{id}}" before init and id when the message is missing.
func Localize(id string, data ...map[string]any) string {
	ensureCache()
	return lookupMessage(cache.local, id, data...)
}

// Unlocalize is Localize in English, used for logs.
func Unlocalize(id string, data ...map[string]any) string {
	ensureCache()
	return lookupMessage(cache.localcatalog, id, data...)
}

func lookupMessage(t *i18n.Localizer, id string, data ...map[string]any) string {
	if t == nil {
		return fmt.Sprintf("{{%s}}", id)
	}

	cfg := makeConfig(t, id, data...)
	if localized, err := t.Localize(cfg); err == nil {
		return localized
	}

	return id
}

func makeConfig(t *i18n.Localizer, id string, data ...map[string]any) *i18n.LocalizeConfig {
	return &i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: makeTemplate(t, data...),
	}
}

func makeTemplate(t *i18n.Localizer, data ...map[string]any) map[string]any {
	if len(data) <= 0 {
		return nil
	}

	template := make(map[string]any)

	for key, value := range data[0] {
		switch v := value.(type) {
		case string:
			template[key] = lookupMessage(t, v)
		case int:
			template[key] = strconv.Itoa(v)
		default:
			template[key] = v
		}
	}

	return template
}

func ensureCache() {
	if cache == nil {
		cache = New()
	}
}
