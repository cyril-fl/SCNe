package classes

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

func TestClassCacheSupportsRuntimeInsertAndCachedHit(t *testing.T) {
	originalPath, originalCache := filePath, cachedConfig
	t.Cleanup(func() { filePath, cachedConfig = originalPath, originalCache })

	path := filepath.Join(t.TempDir(), "classes.yaml")
	contents := "- name: mock.class\n  playable: true\n  bonus: []\n  armor_set: mock.set\n"
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	filePath = path
	cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

	if values, err := ListWithError(); err != nil || len(values) != 1 {
		t.Fatalf("ListWithError() = %#v, %v", values, err)
	}
	if _, err := InsertIfIsnt("mock.runtime.class"); err != nil {
		t.Fatalf("InsertIfIsnt() error = %v", err)
	}
	if _, err := InsertIfIsnt("mock.runtime.class"); err != nil {
		t.Fatalf("second InsertIfIsnt() error = %v", err)
	}
	if len(List()) != 2 {
		t.Fatalf("List() length = %d, want 2", len(List()))
	}
}
