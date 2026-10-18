package oracle

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

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
