package characters

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

func TestLoadFromSourceLoadsRootAndNestedCharacters(t *testing.T) {
	originalFolder := folderConfigPath
	originalConfig := cachedConfig
	t.Cleanup(func() {
		folderConfigPath = originalFolder
		cachedConfig = originalConfig
		logs.Init(nil)
	})
	logs.Init(cachingTestLogRegistry{})

	folder := t.TempDir()
	writeCharacterConfig(t, filepath.Join(folder, "root.yaml"), "mock.character.root")
	nested := filepath.Join(folder, "npcs")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	writeCharacterConfig(t, filepath.Join(nested, "nested.yaml"), "mock.character.nested")

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
		t.Fatalf("cached character count = %d, want %d", got, want)
	}
	for _, name := range []string{"mock.character.root", "mock.character.nested"} {
		if !IsLocked(name) {
			t.Errorf("IsLocked(%q) = false, want true", name)
		}
	}
}

func writeCharacterConfig(t *testing.T, path, name string) {
	t.Helper()
	contents := "name: " + name + "\n" +
		"description: mock.description\n" +
		"rarity: mock.rarity\n" +
		"variety: mock.variety\n" +
		"class: mock-class\n" +
		"race: mock-race\n" +
		"stats: {}\n" +
		"equipment: {}\n" +
		"inventory: []\n" +
		"wallet: {}\n"
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadFromFileRegistersNonEmptyAttributes(t *testing.T) {
	folder := t.TempDir()
	loaded := filepath.Join(folder, "loaded.yaml")
	if err := os.WriteFile(loaded, []byte("name: mock.character\nrarity: mock.character.rarity\nvariety: mock.character.variety\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	empty := filepath.Join(folder, "empty.yaml")
	if err := os.WriteFile(empty, []byte("name: mock.character.empty\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := loadFromFile(loaded); err != nil {
		t.Fatalf("loadFromFile(loaded) error = %v", err)
	}
	if _, err := loadFromFile(empty); err != nil {
		t.Fatalf("loadFromFile(empty) error = %v", err)
	}

	if !rarity.Assert("mock.character.rarity") || !variety.Assert("mock.character.variety") {
		t.Error("loaded rarity and variety should be registered")
	}
	if rarity.Assert("") || variety.Assert("") {
		t.Error("empty rarity and variety should not be registered")
	}
}
