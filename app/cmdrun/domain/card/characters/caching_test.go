package characters

import (
	"os"
	"path/filepath"
	"testing"

	"solopg/app/shared/services/logs"
)

type cachingTestLogRegistry struct{}

func (cachingTestLogRegistry) Register(*logs.Log) error { return nil }

func TestLoadFromSourceLoadsRootAndNestedCharacters(t *testing.T) {
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
	writeCharacterConfig(t, filepath.Join(folder, "root.yaml"), "mock.character.root")
	nested := filepath.Join(folder, "npcs")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	writeCharacterConfig(t, filepath.Join(nested, "nested.yaml"), "mock.character.nested")

	folderConfigPath = folder
	cachedConfig = nil
	lockedNames = nil

	if err := loadFromSource(); err != nil {
		t.Fatalf("loadFromSource() error = %v", err)
	}

	if got, want := len(cachedConfig), 2; got != want {
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
		"description: A test character\n" +
		"rarity: A\n" +
		"variety: character_card\n" +
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
