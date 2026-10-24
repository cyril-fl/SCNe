package hint

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"solopg/app/shared/services/logs"
)

type cachingTestLogRegistry struct{}

func (cachingTestLogRegistry) Register(*logs.Log) error { return nil }

func TestHintCacheLoadsRecursivelyAndReturnsValueCopies(t *testing.T) {
	originalPath, originalCache := folderConfigPath, cachedConfig
	t.Cleanup(func() { folderConfigPath, cachedConfig = originalPath, originalCache })

	folder := t.TempDir()
	nested := filepath.Join(folder, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folder, "root.yaml"), []byte("name: mock.hint.root\nvalues: [mock.one]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "nested.yaml"), []byte("name: mock.hint.nested\nvalues: [mock.two]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	folderConfigPath = folder
	cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

	values, err := ListWithError()
	if err != nil || len(values) != 2 {
		t.Fatalf("ListWithError() = %#v, %v", values, err)
	}
	inner, err := GetByName("mock.hint.root")
	if err != nil {
		t.Fatal(err)
	}
	copied := inner.Values()
	copied[0] = "mock.changed"
	if inner.Values()[0] != "mock.one" {
		t.Fatal("hint values were not copied")
	}
}

func TestRollWithoutValuesReturnsNil(t *testing.T) {
	if got := Roll(newFromYamlConfig(yamlConfig{Name: "mock.hint.empty"})); got != nil {
		t.Fatalf("Roll() = %#v, want nil", got)
	}
}

func TestLoadFromFilesKeepsValidConfigsAndCollectsErrors(t *testing.T) {
	t.Cleanup(func() { logs.Init(nil) })
	logs.Init(cachingTestLogRegistry{})

	folder := t.TempDir()
	valid := filepath.Join(folder, "valid.yaml")
	if err := os.WriteFile(valid, []byte("name: mock.valid\nvalues: [mock.one]\n"), 0o644); err != nil {
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
