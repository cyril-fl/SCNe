package variety

import (
	"os"
	"path/filepath"
	"testing"

	"solopg/app/shared/types/primitive"
)

func TestMakeDefaultUpdatesCacheWithoutDuplicates(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() { fileConfigPath, cachedConfig = originalPath, originalCache })

	path := filepath.Join(t.TempDir(), "variety.yaml")
	if err := os.WriteFile(path, []byte("values: [mock.variety]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fileConfigPath = path
	cachedConfig = primitive.NewCache[yamlConfig]()

	MakeDefault("mock.runtime")
	MakeDefault("mock.runtime")
	values, err := ListWithError()
	if err != nil || len(values) != 2 || values[1] != "mock.runtime" {
		t.Fatalf("values = %#v, %v", values, err)
	}
}
