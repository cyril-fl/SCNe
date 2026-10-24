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
		config.Configs = []yamlConfig{{ID: "race.template", Name: "race.template"}}
		config.LockedIDs = map[string]struct{}{"race.template": {}}
	})

	if !IsLocked("race.template") {
		t.Fatal("expected the template race to be locked")
	}

	if err := AddInConfig(Template{ID: "custom", Name: "race.custom"}); err != nil {
		t.Fatalf("adding a custom race failed: %v", err)
	}
	if IsLocked("custom") {
		t.Fatal("expected a race added at runtime not to be locked")
	}
}

func TestReplaceInConfigReplacesByIDOrAdds(t *testing.T) {
	originalConfig := cachedConfig
	t.Cleanup(func() {
		cachedConfig = originalConfig
	})

	cachedConfig = plumprimitive.NewCache[cacheValue]()
	cachedConfig.Update(func(config *cacheValue) {
		config.Configs = []yamlConfig{{ID: "old", Name: "race.old"}, {ID: "other", Name: "race.other"}}
	})

	if err := ReplaceInConfig("old", Template{ID: "old", Name: "race.new"}); err != nil {
		t.Fatalf("replacing a race failed: %v", err)
	}
	if got := NameOf("old"); got != "race.new" {
		t.Fatalf("NameOf(old) = %q, want race.new", got)
	}
	if len(List()) != 2 {
		t.Fatalf("races = %v, want 2 races", List())
	}

	if err := ReplaceInConfig("missing", Template{ID: "added", Name: "race.added"}); err != nil {
		t.Fatalf("adding a missing race failed: %v", err)
	}
	if FindByID("added") == nil || len(List()) != 3 {
		t.Fatalf("races = %v, want race.added appended", List())
	}
}

func TestAssertAndNameOfUseIDs(t *testing.T) {
	originalConfig := cachedConfig
	t.Cleanup(func() {
		cachedConfig = originalConfig
	})

	cachedConfig = plumprimitive.NewCache[cacheValue]()
	cachedConfig.Update(func(config *cacheValue) {
		config.Configs = []yamlConfig{{ID: "known", Name: "race.known"}}
	})

	if !Assert("known") {
		t.Fatal("Assert rejected a known race ID")
	}
	if Assert("race.known") || Assert("unknown") || Assert("") {
		t.Fatal("Assert accepted a name, an unknown ID or an empty ID")
	}
	if got := NameOf("unknown"); got != "unknown" {
		t.Fatalf("NameOf(unknown) = %q, want the ID back", got)
	}
}
