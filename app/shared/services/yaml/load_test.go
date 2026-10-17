package yaml

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"solopg/app/shared/services/logs"
)

type loadTestLogRegistrar struct{}

func (loadTestLogRegistrar) Register(*logs.Log) error { return nil }

func TestMain(m *testing.M) {
	logs.Init(loadTestLogRegistrar{})
	code := m.Run()
	logs.Init(nil)
	os.Exit(code)
}

func TestLoadFromFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "config.yaml")
	writeLoadTestFile(t, file, "name: campaign\ncount: 3\n")

	type config struct {
		Name  string `yaml:"name"`
		Count int    `yaml:"count"`
	}

	got, err := LoadFromFile[config](file)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}
	if want := (&config{Name: "campaign", Count: 3}); !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadFromFile() = %#v, want %#v", got, want)
	}
}

func TestLoadListFromFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "items.yaml")
	writeLoadTestFile(t, file, "- alpha\n- beta\n")

	got, err := LoadListFromFile[string](file)
	if err != nil {
		t.Fatalf("LoadListFromFile() error = %v", err)
	}
	if want := []string{"alpha", "beta"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadListFromFile() = %#v, want %#v", got, want)
	}
}

func TestLoadFromFileReturnsReadAndUnmarshalErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.yaml")
	if _, err := LoadFromFile[map[string]string](missing); err == nil || !strings.Contains(err.Error(), missing) {
		t.Fatalf("LoadFromFile() error = %v, want an error containing %q", err, missing)
	}

	file := filepath.Join(t.TempDir(), "invalid.yaml")
	writeLoadTestFile(t, file, "name: [invalid\n")
	if _, err := LoadFromFile[map[string]string](file); err == nil || !strings.Contains(err.Error(), file) {
		t.Fatalf("LoadFromFile() error = %v, want an error containing %q", err, file)
	}
}

func TestGetFolderEntries(t *testing.T) {
	folder := t.TempDir()
	writeLoadTestFile(t, filepath.Join(folder, "b.yaml"), "value: b\n")
	writeLoadTestFile(t, filepath.Join(folder, "a.yaml"), "value: a\n")
	if err := os.Mkdir(filepath.Join(folder, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}

	files, err := GetFolderFiles(folder)
	if err != nil {
		t.Fatalf("GetFolderFiles() error = %v", err)
	}
	if want := []string{"a.yaml", "b.yaml"}; !reflect.DeepEqual(files, want) {
		t.Fatalf("GetFolderFiles() = %#v, want %#v", files, want)
	}

	directories, err := GetFolderDirectories(folder)
	if err != nil {
		t.Fatalf("GetFolderDirectories() error = %v", err)
	}
	if want := []string{"nested"}; !reflect.DeepEqual(directories, want) {
		t.Fatalf("GetFolderDirectories() = %#v, want %#v", directories, want)
	}
}

func TestGetFilesFromSource(t *testing.T) {
	folder := t.TempDir()
	rootFile := filepath.Join(folder, "root.yaml")
	nested := filepath.Join(folder, "nested")
	nestedFile := filepath.Join(nested, "nested.yaml")
	writeLoadTestFile(t, rootFile, "root: true\n")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	writeLoadTestFile(t, nestedFile, "nested: true\n")

	nonRecursive, err := GetFilesFromSource(folder, false)
	if err != nil {
		t.Fatalf("GetFilesFromSource(false) error = %v", err)
	}
	if want := []string{rootFile}; !reflect.DeepEqual(nonRecursive, want) {
		t.Fatalf("GetFilesFromSource(false) = %#v, want %#v", nonRecursive, want)
	}

	recursive, err := GetFilesFromSource(folder, true)
	if err != nil {
		t.Fatalf("GetFilesFromSource(true) error = %v", err)
	}
	if want := []string{nestedFile, rootFile}; !reflect.DeepEqual(recursive, want) {
		t.Fatalf("GetFilesFromSource(true) = %#v, want %#v", recursive, want)
	}

	single, err := GetFilesFromSource(rootFile, true)
	if err != nil {
		t.Fatalf("GetFilesFromSource(file) error = %v", err)
	}
	if want := []string{rootFile}; !reflect.DeepEqual(single, want) {
		t.Fatalf("GetFilesFromSource(file) = %#v, want %#v", single, want)
	}
}

func TestGetFilesFromSourceReturnsErrorForMissingSource(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	if _, err := GetFilesFromSource(missing, false); err == nil || !strings.Contains(err.Error(), missing) {
		t.Fatalf("GetFilesFromSource() error = %v, want an error containing %q", err, missing)
	}
}

func writeLoadTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
