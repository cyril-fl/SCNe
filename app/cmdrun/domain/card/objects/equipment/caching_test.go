package equipment

import (
	"os"
	"path/filepath"
	"testing"

	"scne/app/cmdrun/domain/card/attributes/objectcategory"
	"scne/app/cmdrun/domain/card/attributes/rarity"
	"scne/app/cmdrun/domain/card/attributes/variety"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"scne/app/shared/services/logs"
)

type cachingTestLogRegistry struct{}

func (cachingTestLogRegistry) Register(*logs.Log) error { return nil }

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

func TestLoadFromFilesKeepsValidConfigsAndCollectsErrors(t *testing.T) {
	t.Cleanup(func() { logs.Init(nil) })
	logs.Init(cachingTestLogRegistry{})

	folder := t.TempDir()
	valid := filepath.Join(folder, "valid.yaml")
	if err := os.WriteFile(valid, []byte("name: mock.valid\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	invalid := filepath.Join(folder, "invalid.yaml")
	if err := os.WriteFile(invalid, []byte("name: [invalid\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	configs, errs := loadFromFiles([]string{valid, invalid})
	if len(configs) != 1 || configs[0].Name != "mock.valid" {
		t.Errorf("configs = %#v, want only mock.valid", configs)
	}
	if len(errs) != 1 {
		t.Errorf("errs = %v, want 1 error", errs)
	}
}
