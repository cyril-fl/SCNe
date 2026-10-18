package objectcategory

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

func TestListUsesLoadedValueAfterSourceRemoval(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() { fileConfigPath, cachedConfig = originalPath, originalCache })

	path := filepath.Join(t.TempDir(), "categories.yaml")
	if err := os.WriteFile(path, []byte("values: [mock.category]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fileConfigPath = path
	cachedConfig = plumprimitive.NewCache[yamlConfig]()

	values, err := ListWithError()
	if err != nil || len(values) != 1 || values[0] != "mock.category" {
		t.Fatalf("ListWithError() = %#v, %v", values, err)
	}
	values[0] = "mock.changed"
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	values, err = ListWithError()
	if err != nil || len(values) != 1 || values[0] != "mock.category" {
		t.Fatalf("cached ListWithError() = %#v, %v", values, err)
	}
}
