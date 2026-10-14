package set

import "testing"

func TestSetAddHasRemoveAndSize(t *testing.T) {
	s := New[string]()

	if s.Size() != 0 {
		t.Fatalf("new set size = %d, want 0", s.Size())
	}
	if s.Has("one") {
		t.Fatal("new set unexpectedly contains one")
	}

	s.Add("one")
	s.Add("two")
	s.Add("one")

	if !s.Has("one") || !s.Has("two") {
		t.Fatal("set does not contain all added values")
	}
	if s.Size() != 2 {
		t.Fatalf("set size = %d, want 2", s.Size())
	}

	s.Remove("one")
	if s.Has("one") {
		t.Fatal("removed value is still in set")
	}
	if s.Size() != 1 {
		t.Fatalf("set size after removal = %d, want 1", s.Size())
	}
}

func TestSetToSlice(t *testing.T) {
	s := New[string]()
	s.Add("one")
	s.Add("two")

	values := s.ToSlice()
	if len(values) != 2 {
		t.Fatalf("slice length = %d, want 2", len(values))
	}

	for _, value := range []string{"one", "two"} {
		found := false
		for _, got := range values {
			if got == value {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("slice does not contain %q: %v", value, values)
		}
	}
}
