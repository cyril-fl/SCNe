package rarity

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"scne/app/shared/services/logs"
)

type cachingTestLogRegistry struct{}

func (cachingTestLogRegistry) Register(*logs.Log) error { return nil }

func TestDefaultAndListWithErrorUseConfiguredValues(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() { fileConfigPath, cachedConfig = originalPath, originalCache })

	path := filepath.Join(t.TempDir(), "rarity.yaml")
	if err := os.WriteFile(path, []byte("values: [mock.rarity]\ndefault: mock.rarity\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fileConfigPath = path
	cachedConfig = plumprimitive.NewCache[yamlConfig]()

	if got := Default(); got != "mock.rarity" {
		t.Fatalf("Default() = %q", got)
	}
	values, err := ListWithError()
	if err != nil || len(values) != 1 || values[0] != "mock.rarity" {
		t.Fatalf("ListWithError() = %#v, %v", values, err)
	}
}

func TestRarityIsOptional(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() { fileConfigPath, cachedConfig = originalPath, originalCache })

	for _, path := range []string{"", filepath.Join(t.TempDir(), "missing.yaml")} {
		fileConfigPath = path
		cachedConfig = plumprimitive.NewCache[yamlConfig]()

		values, err := ListWithError()
		if err != nil || len(values) != 0 {
			t.Fatalf("ListWithError() with path %q = %#v, %v; want no values, no error", path, values, err)
		}
		if got := Default(); got != "" || !got.Validate() {
			t.Fatalf("Default() with path %q = %q; want valid empty rarity", path, got)
		}
	}
}

func TestInvalidRarityFileStillFails(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() {
		fileConfigPath, cachedConfig = originalPath, originalCache
		logs.Init(nil)
	})
	logs.Init(cachingTestLogRegistry{})

	path := filepath.Join(t.TempDir(), "rarity.yaml")
	if err := os.WriteFile(path, []byte("values: [invalid\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fileConfigPath = path
	cachedConfig = plumprimitive.NewCache[yamlConfig]()

	if _, err := ListWithError(); err == nil {
		t.Fatal("expected invalid YAML to fail loading")
	}
}
