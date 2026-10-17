package locations

import (
	"os"
	"path/filepath"
	"testing"

	"solopg/app/shared/services/logs"
	"solopg/app/shared/types/primitive"
)

type cachingTestLogRegistry struct{}

func (cachingTestLogRegistry) Register(*logs.Log) error { return nil }

func TestLoadFromSourceLoadsRootAndNestedLocations(t *testing.T) {
	originalFolder := folderConfigPath
	originalConfig := cachedConfig
	t.Cleanup(func() {
		folderConfigPath = originalFolder
		cachedConfig = originalConfig
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
	cachedConfig = primitive.NewCache[cacheValue]()

	if _, err := cachedConfig.Load(loadFromSource); err != nil {
		t.Fatalf("loadFromSource() error = %v", err)
	}

	config, loaded := cachedConfig.Value()
	if !loaded {
		t.Fatal("cache should be loaded")
	}
	if got, want := len(config.Configs), 2; got != want {
		t.Fatalf("cached location count = %d, want %d", got, want)
	}
	for _, name := range []string{"mock.location.root", "mock.location.nested"} {
		if !IsLocked(name) {
			t.Errorf("IsLocked(%q) = false, want true", name)
		}
	}
}

func TestLoadFromSourceDoesNotPublishPartialConfiguration(t *testing.T) {
	originalFolder := folderConfigPath
	originalConfig := cachedConfig
	t.Cleanup(func() {
		folderConfigPath = originalFolder
		cachedConfig = originalConfig
		logs.Init(nil)
	})
	logs.Init(cachingTestLogRegistry{})

	folder := t.TempDir()
	writeLocationConfig(t, filepath.Join(folder, "valid.yaml"), "mock.location.valid")
	if err := os.WriteFile(filepath.Join(folder, "invalid.yaml"), []byte("name: [invalid\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	folderConfigPath = folder
	cachedConfig = primitive.NewCache[cacheValue]()

	if _, err := ListWithError(); err == nil {
		t.Fatal("expected invalid YAML to fail loading")
	}
	value, loaded := cachedConfig.Value()
	if !loaded || len(value.Configs) != 0 || len(value.LockedNames) != 0 {
		t.Fatalf("partial cache published: loaded=%t value=%#v", loaded, value)
	}
}

func writeLocationConfig(t *testing.T, path, name string) {
	t.Helper()
	contents := "name: " + name + "\n" +
		"description: " + name + ".description\n" +
		"rarity: mock.rarity\n" +
		"variety: mock.variety\n" +
		"effects: []\n"
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
