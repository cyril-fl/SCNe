package locations

import (
	"os"
	"path/filepath"
	"testing"

	"solopg/app/shared/services/logs"
)

type cachingTestLogRegistry struct{}

func (cachingTestLogRegistry) Register(*logs.Log) error { return nil }

func TestLoadFromSourceLoadsRootAndNestedLocations(t *testing.T) {
	originalFolder := folderConfigPath
	originalConfig := cachedConfig
	originalLockedNames := lockedNames
	t.Cleanup(func() {
		folderConfigPath = originalFolder
		cachedConfig = originalConfig
		lockedNames = originalLockedNames
		logs.Init(nil)
	})
	logs.Init(cachingTestLogRegistry{})

	folder := t.TempDir()
	writeLocationConfig(t, filepath.Join(folder, "root.yaml"), "mock.location.root")
	nested := filepath.Join(folder, "special")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	writeLocationConfig(t, filepath.Join(nested, "nested.yaml"), "mock.location.nested")

	folderConfigPath = folder
	cachedConfig = nil
	lockedNames = nil

	if err := loadFromSource(); err != nil {
		t.Fatalf("loadFromSource() error = %v", err)
	}

	if got, want := len(cachedConfig), 2; got != want {
		t.Fatalf("cached location count = %d, want %d", got, want)
	}
	for _, name := range []string{"mock.location.root", "mock.location.nested"} {
		if !IsLocked(name) {
			t.Errorf("IsLocked(%q) = false, want true", name)
		}
	}
}

func writeLocationConfig(t *testing.T, path, name string) {
	t.Helper()
	contents := "name: " + name + "\n" +
		"description: " + name + ".description\n" +
		"rarity: E\n" +
		"variety: location_card\n" +
		"effects: []\n"
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
