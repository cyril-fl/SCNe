package stats

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

func TestBasicStatsAndNamesAreCopiedFromCache(t *testing.T) {
	originalPath, originalCache := fileConfigPath, cachedConfig
	t.Cleanup(func() { fileConfigPath, cachedConfig = originalPath, originalCache })

	path := filepath.Join(t.TempDir(), "stats.yaml")
	contents := "names: [mock.stat]\nbaseStats:\n  mock.stat: 7\n"
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	fileConfigPath = path
	cachedConfig = plumprimitive.NewCache[yamlConfig]()

	names, err := ListWithError()
	if err != nil || len(names) != 1 || names[0] != "mock.stat" {
		t.Fatalf("ListWithError() = %#v, %v", names, err)
	}
	base, err := GetBasicWithError()
	if err != nil || base[Stat("mock.stat")] != 7 {
		t.Fatalf("GetBasicWithError() = %#v, %v", base, err)
	}
	base[Stat("mock.stat")] = 99
	base, err = GetBasicWithError()
	if err != nil || base[Stat("mock.stat")] != 7 {
		t.Fatalf("cached basic stats = %#v, %v", base, err)
	}
}
