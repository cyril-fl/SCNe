package hint

import (
	"os"
	"path/filepath"
	"testing"

	"solopg/app/shared/types/primitive"
)

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
	cachedConfig = primitive.NewCache[[]yamlConfig]()

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
