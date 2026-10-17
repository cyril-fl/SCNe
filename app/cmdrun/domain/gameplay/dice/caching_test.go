package dice

import (
	"os"
	"path/filepath"
	"testing"

	"solopg/app/shared/types/primitive"
)

func TestDiceCacheLoadsEmptyAndKeepsItLoaded(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() { fileConfigPath, cachedConfig = originalPath, originalCache })

	path := filepath.Join(t.TempDir(), "dice.yaml")
	if err := os.WriteFile(path, []byte("[]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fileConfigPath = path
	cachedConfig = primitive.NewCache[[]yamlConfig]()

	values, err := ListWithError()
	if err != nil || len(values) != 0 || !cachedConfig.IsLoaded() {
		t.Fatalf("ListWithError() = %#v, %v", values, err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	values, err = ListWithError()
	if err != nil || len(values) != 0 {
		t.Fatalf("cached ListWithError() = %#v, %v", values, err)
	}
}
