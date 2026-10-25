package oracle

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"scne/app/shared/services/logs"
)

type cachingTestLogRegistry struct{}

func (cachingTestLogRegistry) Register(*logs.Log) error { return nil }

func TestOracleCacheLoadsValidTemporaryYAML(t *testing.T) {
	originalPath, originalCache := folderConfigPath, cachedConfig
	t.Cleanup(func() { folderConfigPath, cachedConfig = originalPath, originalCache })

	folder := t.TempDir()
	path := filepath.Join(folder, "oracle.yaml")
	contents := "name: mock.oracle\ndice: 1\ntype: mock.type\nvisible: true\nintervals:\n  - min: 1\n    max: 1\n    critical: false\n    result: mock.result\n"
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	folderConfigPath = folder
	cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

	values, err := ListWithError()
	if err != nil || len(values) != 1 || values[0].name != "mock.oracle" {
		t.Fatalf("ListWithError() = %#v, %v", values, err)
	}
}

func TestLoadFromFilesKeepsValidConfigsAndCollectsErrors(t *testing.T) {
	t.Cleanup(func() { logs.Init(nil) })
	logs.Init(cachingTestLogRegistry{})

	folder := t.TempDir()
	valid := filepath.Join(folder, "valid.yaml")
	if err := os.WriteFile(valid, []byte("name: mock.valid\ndice: 1\ntype: mock.type\nvisible: true\nintervals:\n  - min: 1\n    max: 1\n    critical: false\n    result: mock.result\n"), 0o644); err != nil {
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
