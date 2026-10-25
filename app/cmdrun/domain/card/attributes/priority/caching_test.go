package priority

import (
	"os"
	"path/filepath"
	"testing"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"scne/app/shared/services/logs"
)

type cachingTestLogRegistry struct{}

func (cachingTestLogRegistry) Register(*logs.Log) error { return nil }

func TestPriorityCacheSupportsLookupAndNameUniqueness(t *testing.T) {
	originalPath, originalCache := filePath, cachedConfig
	t.Cleanup(func() {
		filePath, cachedConfig = originalPath, originalCache
		logs.Init(nil)
	})
	logs.Init(cachingTestLogRegistry{})

	path := filepath.Join(t.TempDir(), "priorities.yaml")
	contents := "- name: mock.priority.low\n  weight: 1\n"
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	filePath = path
	cachedConfig = plumprimitive.NewCache[[]yamlConfig]()

	values, err := ListWithError()
	if err != nil || len(values) != 1 {
		t.Fatalf("ListWithError() = %#v, %v", values, err)
	}
	if FindByName("mock.priority.low") == nil || len(ListNames()) != 1 {
		t.Fatal("priority lookup by name failed")
	}
	if err := AddInConfig(Template{Name: "mock.priority.high", Weight: 2}); err != nil {
		t.Fatalf("AddInConfig() error = %v", err)
	}
	if err := AddInConfig(Template{Name: "mock.priority.high", Weight: 3}); err == nil {
		t.Fatal("duplicate priority name should fail")
	}
	if len(List()) != 2 {
		t.Fatalf("List() length = %d, want 2", len(List()))
	}
}
