package editor

import (
	"slices"
	"testing"
)

func TestEditorPrefersVisualThenEditorThenVi(t *testing.T) {
	t.Setenv("VISUAL", "code --wait")
	t.Setenv("EDITOR", "nano")
	if got := command(); !slices.Equal(got, []string{"code", "--wait"}) {
		t.Fatalf("command() = %v, want VISUAL split", got)
	}

	t.Setenv("VISUAL", "")
	if got := command(); !slices.Equal(got, []string{"nano"}) {
		t.Fatalf("command() = %v, want EDITOR", got)
	}

	t.Setenv("EDITOR", "")
	if got := command(); !slices.Equal(got, []string{"vi"}) {
		t.Fatalf("command() = %v, want vi", got)
	}
}
