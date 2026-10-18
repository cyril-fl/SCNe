package rarity

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

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
