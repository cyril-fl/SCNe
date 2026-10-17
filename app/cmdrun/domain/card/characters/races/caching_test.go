package races

import (
	"testing"

	"solopg/app/shared/types/primitive"
)

func TestIsLockedTracksTemplateEntriesOnly(t *testing.T) {
	originalConfig := cachedConfig
	t.Cleanup(func() {
		cachedConfig = originalConfig
	})

	cachedConfig = primitive.NewCache[cacheValue]()
	cachedConfig.Update(func(config *cacheValue) {
		config.Configs = []yamlConfig{{Name: "race.template"}}
		config.LockedNames = map[string]struct{}{"race.template": {}}
	})

	if !IsLocked("race.template") {
		t.Fatal("expected the template race to be locked")
	}

	if err := AddInConfig(Template{Name: "race.custom"}); err != nil {
		t.Fatalf("adding a custom race failed: %v", err)
	}
	if IsLocked("race.custom") {
		t.Fatal("expected a race added at runtime not to be locked")
	}
}
