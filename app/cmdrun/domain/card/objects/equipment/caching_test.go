package equipment

import (
	"os"
	"path/filepath"
	"testing"

	"solopg/app/shared/types/primitive"
)

func TestLoadFromSourcePublishesAllTemporaryFilesAtomically(t *testing.T) {
	originalPath, originalCache := folderConfigPath, cachedConfig
	t.Cleanup(func() { folderConfigPath, cachedConfig = originalPath, originalCache })

	folder := t.TempDir()
	nested := filepath.Join(folder, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		filepath.Join(folder, "root.yaml"),
		filepath.Join(nested, "nested.yaml"),
	} {
		if err := os.WriteFile(path, []byte("name: mock.gear\nset: mock.set\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	folderConfigPath = folder
	cachedConfig = primitive.NewCache[[]yamlConfig]()

	configs, err := cachedConfig.Load(loadFromSource)
	if err != nil || len(configs) != 2 {
		t.Fatalf("loadFromSource() = %#v, %v", configs, err)
	}
	if err := os.Remove(filepath.Join(nested, "nested.yaml")); err != nil {
		t.Fatal(err)
	}
	again, err := cachedConfig.Load(loadFromSource)
	if err != nil || len(again) != 2 {
		t.Fatalf("cached load = %#v, %v", again, err)
	}
}
