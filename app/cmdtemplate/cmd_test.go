package cmdtemplate

import (
	"os"
	"path/filepath"
	"testing"
)

// useEditor sets EDITOR to a script running body with the file path as $1.
func useEditor(t *testing.T, body string) {
	t.Helper()
	script := filepath.Join(t.TempDir(), "editor.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", script)
}

func checkName(data []byte) (string, error) {
	return kinds()[0].Validate(data)
}

func TestEditOnceCancelsWhenUnchanged(t *testing.T) {
	useEditor(t, "true")

	data, _, err := editOnce([]byte("# name: \"\"\n"), checkName, "scne import --class")
	if data != nil || err != nil {
		t.Fatalf("editOnce() = %q, %v, want a cancel", data, err)
	}
}

func TestEditOnceKeepsTheDraftOnError(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	useEditor(t, `echo 'name: ""' > "$1"`)

	if _, _, err := editOnce([]byte("# name: \"\"\n"), checkName, "scne import --class"); err == nil {
		t.Fatal("editOnce() error = nil, want the validation error")
	}

	drafts, _ := filepath.Glob(filepath.Join(tmp, "scne-*.yaml"))
	if len(drafts) != 1 {
		t.Fatalf("drafts = %v, want the edited file kept", drafts)
	}
}

func TestEditOnceRemovesTheFileWhenDone(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	useEditor(t, "true")

	editOnce([]byte("# name: \"\"\n"), checkName, "scne import --class")

	if drafts, _ := filepath.Glob(filepath.Join(tmp, "scne-*.yaml")); len(drafts) != 0 {
		t.Fatalf("drafts = %v, want none after a cancel", drafts)
	}
}

func TestEditOnceReturnsTheValidEntry(t *testing.T) {
	useEditor(t, `echo 'name: ranger' > "$1"`)

	data, name, err := editOnce([]byte("# name: \"\"\n"), checkName, "scne import --class")
	if err != nil || name != "ranger" || string(data) != "name: ranger\n" {
		t.Fatalf("editOnce() = %q, %q, %v", data, name, err)
	}
}
