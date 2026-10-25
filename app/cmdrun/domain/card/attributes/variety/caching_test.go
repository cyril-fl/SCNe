package variety

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"solopg/app/shared/services/logs"
)

type cachingTestLogRegistry struct{}

func (cachingTestLogRegistry) Register(*logs.Log) error { return nil }

func TestMakeDefaultUpdatesCacheWithoutDuplicates(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() { fileConfigPath, cachedConfig = originalPath, originalCache })

	path := filepath.Join(t.TempDir(), "variety.yaml")
	if err := os.WriteFile(path, []byte("values: [mock.variety]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fileConfigPath = path
	cachedConfig = plumprimitive.NewCache[yamlConfig]()

	MakeDefault("mock.runtime")
	MakeDefault("mock.runtime")
	values, err := ListWithError()
	if err != nil || len(values) != 2 || values[1] != "mock.runtime" {
		t.Fatalf("values = %#v, %v", values, err)
	}
}

func TestDefaultUsesConfiguredValue(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() { fileConfigPath, cachedConfig = originalPath, originalCache })

	path := filepath.Join(t.TempDir(), "variety.yaml")
	if err := os.WriteFile(path, []byte("values: [mock.first, mock.variety]\ndefault: mock.variety\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fileConfigPath = path
	cachedConfig = plumprimitive.NewCache[yamlConfig]()

	if got := Default(); got != "mock.variety" {
		t.Fatalf("Default() = %q", got)
	}
}

func TestVarietyIsOptional(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() { fileConfigPath, cachedConfig = originalPath, originalCache })

	for _, path := range []string{"", filepath.Join(t.TempDir(), "missing.yaml")} {
		fileConfigPath = path
		cachedConfig = plumprimitive.NewCache[yamlConfig]()

		values, err := ListWithError()
		if err != nil || len(values) != 0 {
			t.Fatalf("ListWithError() with path %q = %#v, %v; want no values, no error", path, values, err)
		}
		if got := Default(); got != "card" || !got.Validate() {
			t.Fatalf("Default() with path %q = %q; want valid card variety", path, got)
		}
	}
}

func TestInvalidVarietyFileStillFails(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() {
		fileConfigPath, cachedConfig = originalPath, originalCache
		logs.Init(nil)
	})
	logs.Init(cachingTestLogRegistry{})

	path := filepath.Join(t.TempDir(), "variety.yaml")
	if err := os.WriteFile(path, []byte("values: [invalid\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fileConfigPath = path
	cachedConfig = plumprimitive.NewCache[yamlConfig]()

	if _, err := ListWithError(); err == nil {
		t.Fatal("expected invalid YAML to fail loading")
	}
}
