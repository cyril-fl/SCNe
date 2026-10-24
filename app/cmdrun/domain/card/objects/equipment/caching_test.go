package equipment

import (
	"os"
	"path/filepath"
	"testing"

	"solopg/app/cmdrun/domain/card/attributes/objectcategory"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/variety"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

func TestLoadFromSourcePublishesAllTemporaryFilesAtomically(t *testing.T) {
	originalPath, originalCache := folderConfigPath, cachedConfig
	t.Cleanup(func() { folderConfigPath, cachedConfig = originalPath, originalCache })

	folder := t.TempDir()
	nested := filepath.Join(folder, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		filepath.Join(folder, "root.yaml"),
		filepath.Join(nested, "nested.yaml"),
	} {
		if err := os.WriteFile(path, []byte("name: mock.gear\nset: mock.set\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	folderConfigPath = folder
	cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

	configs, err := cachedConfig.Load(loadFromSource)
	if err != nil || len(configs) != 2 {
		t.Fatalf("loadFromSource() = %#v, %v", configs, err)
	}
	if err := os.Remove(filepath.Join(nested, "nested.yaml")); err != nil {
		t.Fatal(err)
	}
	again, err := cachedConfig.Load(loadFromSource)
	if err != nil || len(again) != 2 {
		t.Fatalf("cached load = %#v, %v", again, err)
	}
}

func TestLoadFromFileRegistersNonEmptyAttributes(t *testing.T) {
	folder := t.TempDir()
	loaded := filepath.Join(folder, "loaded.yaml")
	contents := "name: mock.gear\nrarity: mock.gear.rarity\nvariety: mock.gear.variety\ncategory: mock.gear.category\n"
	if err := os.WriteFile(loaded, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	empty := filepath.Join(folder, "empty.yaml")
	if err := os.WriteFile(empty, []byte("name: mock.gear.empty\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := loadFromFile(loaded); err != nil {
		t.Fatalf("loadFromFile(loaded) error = %v", err)
	}
	if _, err := loadFromFile(empty); err != nil {
		t.Fatalf("loadFromFile(empty) error = %v", err)
	}

	if !rarity.Assert("mock.gear.rarity") || !variety.Assert("mock.gear.variety") || !objectcategory.Category("mock.gear.category").Validate() {
		t.Error("loaded rarity, variety and category should be registered")
	}
	if rarity.Assert("") || variety.Assert("") || objectcategory.Category("").Validate() {
		t.Error("empty rarity, variety and category should not be registered")
	}
}
