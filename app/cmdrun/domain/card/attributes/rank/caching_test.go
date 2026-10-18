package rank

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

func TestListDefaultAndRuntimeValueUseOneCache(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() { fileConfigPath, cachedConfig = originalPath, originalCache })

	path := filepath.Join(t.TempDir(), "ranks.yaml")
	if err := os.WriteFile(path, []byte("values: [mock.rank]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fileConfigPath = path
	cachedConfig = plumprimitive.NewCache[yamlConfig]()

	if got := Default(); got != "mock.rank" {
		t.Fatalf("Default() = %q", got)
	}
	if got := MakeDefault("mock.runtime"); got != "mock.runtime" {
		t.Fatalf("MakeDefault() = %q", got)
	}
	values, err := ListWithError()
	if err != nil || len(values) != 2 {
		t.Fatalf("ListWithError() = %#v, %v", values, err)
	}
}
