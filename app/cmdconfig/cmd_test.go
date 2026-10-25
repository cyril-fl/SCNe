package cmdconfig

import (
	"slices"
	"testing"
)

func TestEditorPrefersVisualThenEditorThenVi(t *testing.T) {
	t.Setenv("VISUAL", "code --wait")
	t.Setenv("EDITOR", "nano")
	if got := editor(); !slices.Equal(got, []string{"code", "--wait"}) {
		t.Fatalf("editor() = %v, want VISUAL split", got)
	}

	t.Setenv("VISUAL", "")
	if got := editor(); !slices.Equal(got, []string{"nano"}) {
		t.Fatalf("editor() = %v, want EDITOR", got)
	}

	t.Setenv("EDITOR", "")
	if got := editor(); !slices.Equal(got, []string{"vi"}) {
		t.Fatalf("editor() = %v, want vi", got)
	}
}
