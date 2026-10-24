package races

import (
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

func TestIsLockedTracksTemplateEntriesOnly(t *testing.T) {
	originalConfig := cachedConfig
	t.Cleanup(func() {
		cachedConfig = originalConfig
	})

	cachedConfig = plumprimitive.NewCache[cacheValue]()
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

func TestReplaceInConfigReplacesByNameOrAdds(t *testing.T) {
	originalConfig := cachedConfig
	t.Cleanup(func() {
		cachedConfig = originalConfig
	})

	cachedConfig = plumprimitive.NewCache[cacheValue]()
	cachedConfig.Update(func(config *cacheValue) {
		config.Configs = []yamlConfig{{Name: "race.old"}, {Name: "race.other"}}
	})

	if err := ReplaceInConfig("race.old", Template{Name: "race.new"}); err != nil {
		t.Fatalf("replacing a race failed: %v", err)
	}
	if FindByName("race.old") != nil || FindByName("race.new") == nil {
		t.Fatalf("races = %v, want race.old replaced by race.new", ListNames())
	}
	if len(List()) != 2 {
		t.Fatalf("races = %v, want 2 races", ListNames())
	}

	if err := ReplaceInConfig("race.missing", Template{Name: "race.added"}); err != nil {
		t.Fatalf("adding a missing race failed: %v", err)
	}
	if FindByName("race.added") == nil || len(List()) != 3 {
		t.Fatalf("races = %v, want race.added appended", ListNames())
	}
}
