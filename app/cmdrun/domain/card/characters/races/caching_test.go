package races

import "testing"

func TestIsLockedTracksTemplateEntriesOnly(t *testing.T) {
	originalConfig := cachedConfig
	originalLockedNames := lockedNames
	t.Cleanup(func() {
		cachedConfig = originalConfig
		lockedNames = originalLockedNames
	})

	cachedConfig = []yamlConfig{{Name: "race.template"}}
	lockedNames = map[string]struct{}{"race.template": {}}

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
