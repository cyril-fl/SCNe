package annotations

import "testing"

func TestAnnotations(t *testing.T) {
	a := Set(NODB, "true").Set("other", "")

	if !a.Has(NODB) || a.Get(NODB) != "true" {
		t.Fatalf("Get(%q) = %q, want true", NODB, a.Get(NODB))
	}
	if !a.Has("other") {
		t.Fatal("Has() ignored a key with an empty value")
	}
	if a.Has("missing") || a.Get("missing") != "" {
		t.Fatal("missing key reported as present")
	}
}
