package theme

import "testing"

func TestMergeKeepsSetColorsAndFillsEmptyOnes(t *testing.T) {
	got := Config{Accent: "#000000"}.Merge(Default)

	if got.Accent != "#000000" {
		t.Fatalf("Accent = %q, want %q", got.Accent, "#000000")
	}
	if got.Muted != Default.Muted {
		t.Fatalf("Muted = %q, want default %q", got.Muted, Default.Muted)
	}
}
