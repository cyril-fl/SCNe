package locations

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/shared/services/logs"
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
	cachedConfig = plumprimitive.NewCache[cacheValue]()

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
	cachedConfig = plumprimitive.NewCache[cacheValue]()

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

func TestLoadFromFileRegistersNonEmptyAttributes(t *testing.T) {
	folder := t.TempDir()
	loaded := filepath.Join(folder, "loaded.yaml")
	if err := os.WriteFile(loaded, []byte("name: mock.location\nrarity: mock.location.rarity\nvariety: mock.location.variety\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	empty := filepath.Join(folder, "empty.yaml")
	if err := os.WriteFile(empty, []byte("name: mock.location.empty\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := loadFromFile(loaded); err != nil {
		t.Fatalf("loadFromFile(loaded) error = %v", err)
	}
	if _, err := loadFromFile(empty); err != nil {
		t.Fatalf("loadFromFile(empty) error = %v", err)
	}

	if !rarity.Assert("mock.location.rarity") || !variety.Assert("mock.location.variety") {
		t.Error("loaded rarity and variety should be registered")
	}
	if rarity.Assert("") || variety.Assert("") {
		t.Error("empty rarity and variety should not be registered")
	}
}
