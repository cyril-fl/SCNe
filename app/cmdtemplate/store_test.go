package cmdtemplate

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestListStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "classes.yaml")
	if err := os.WriteFile(path, []byte("# Playable classes\n- name: mage\n  playable: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := listStore{path: path}

	if err := s.Put("rogue", []byte("name: rogue\nplayable: false\n")); err != nil {
		t.Fatalf("Put(new) error = %v", err)
	}
	if err := s.Put("mage", []byte("name: wizard\nplayable: true\n")); err != nil {
		t.Fatalf("Put(existing) error = %v", err)
	}
	if names, _ := s.Names(); !slices.Equal(names, []string{"wizard", "rogue"}) {
		t.Fatalf("Names() = %v, want the replaced entry in place, then the added one", names)
	}
	if data, _ := os.ReadFile(path); !strings.Contains(string(data), "# Playable classes") {
		t.Fatalf("file lost its comment:\n%s", data)
	}

	if data, err := s.Get("rogue"); err != nil || string(data) != "name: rogue\nplayable: false\n" {
		t.Fatalf("Get() = %q, %v", data, err)
	}
	if err := s.Delete("rogue"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := s.Get("rogue"); !errors.Is(err, errNotFound) {
		t.Fatalf("Get(deleted) error = %v, want errNotFound", err)
	}
}

func TestListStoreCreatesMissingFile(t *testing.T) {
	s := listStore{path: filepath.Join(t.TempDir(), "new", "races.yaml")}

	if err := s.Put("orc", []byte("name: orc\n")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if names, _ := s.Names(); !slices.Equal(names, []string{"orc"}) {
		t.Fatalf("Names() = %v, want [orc]", names)
	}
}

func TestFolderStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "npcs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "npcs", "hero.yaml"), []byte("name: Hero\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := folderStore{dir: dir}

	if err := s.Put("Black smith", []byte("name: Black smith\n")); err != nil {
		t.Fatalf("Put(new) error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "black_smith.yaml")); err != nil {
		t.Fatalf("expected a file named after the entry: %v", err)
	}

	if err := s.Put("Hero", []byte("name: Hero\ndescription: edited\n")); err != nil {
		t.Fatalf("Put(existing) error = %v", err)
	}
	if data, _ := os.ReadFile(filepath.Join(dir, "npcs", "hero.yaml")); !strings.Contains(string(data), "edited") {
		t.Fatalf("expected the existing file to be replaced in its subfolder, got %q", data)
	}

	if names, _ := s.Names(); !slices.Equal(names, []string{"Black smith", "Hero"}) {
		t.Fatalf("Names() = %v", names)
	}
	if err := s.Delete("Hero"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if err := s.Delete("Hero"); !errors.Is(err, errNotFound) {
		t.Fatalf("Delete(deleted) error = %v, want errNotFound", err)
	}
}

func TestFolderStoreMissingFolderIsEmpty(t *testing.T) {
	names, err := folderStore{dir: filepath.Join(t.TempDir(), "missing")}.Names()
	if err != nil || len(names) != 0 {
		t.Fatalf("Names() = %v, %v, want no entry", names, err)
	}
}

func TestReadEntriesSplitsLists(t *testing.T) {
	dir := t.TempDir()
	list := filepath.Join(dir, "list.yaml")
	single := filepath.Join(dir, "single.yaml")
	os.WriteFile(list, []byte("- name: a\n- name: b\n"), 0o644)
	os.WriteFile(single, []byte("name: c\n"), 0o644)

	if entries, err := readEntries(list); err != nil || len(entries) != 2 || string(entries[1]) != "name: b\n" {
		t.Fatalf("readEntries(list) = %q, %v", entries, err)
	}
	if entries, err := readEntries(single); err != nil || len(entries) != 1 {
		t.Fatalf("readEntries(single) = %q, %v", entries, err)
	}
}

func TestErrorLinesAreRemovedBeforeValidation(t *testing.T) {
	data := append(errorLines(errors.New("bad\nworse")), []byte("# kept\nname: a\n")...)

	if got := string(withoutErrorLines(data)); got != "# kept\nname: a\n" {
		t.Fatalf("withoutErrorLines() = %q", got)
	}
}

func TestAskReopenQuitsOnlyOnQ(t *testing.T) {
	tests := map[string]bool{"\n": true, "r\n": true, "q\n": false, "Q\n": false, "": false}
	for input, want := range tests {
		if got := askReopen(strings.NewReader(input), io.Discard); got != want {
			t.Errorf("askReopen(%q) = %v, want %v", input, got, want)
		}
	}
}
