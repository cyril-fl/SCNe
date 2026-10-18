package slot

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

func TestListWithErrorMemoizesAnEmptyConfiguration(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() { fileConfigPath, cachedConfig = originalPath, originalCache })

	path := filepath.Join(t.TempDir(), "slots.yaml")
	if err := os.WriteFile(path, []byte("values: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fileConfigPath = path
	cachedConfig = plumprimitive.NewCache[yamlConfig]()

	values, err := ListWithError()
	if err != nil || len(values) != 0 || !cachedConfig.IsLoaded() {
		t.Fatalf("ListWithError() = %#v, %v; loaded=%t", values, err, cachedConfig.IsLoaded())
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	values, err = ListWithError()
	if err != nil || len(values) != 0 {
		t.Fatalf("cached empty ListWithError() = %#v, %v", values, err)
	}
}
