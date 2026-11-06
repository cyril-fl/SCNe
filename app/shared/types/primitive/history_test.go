package primitive

import (
	"slices"
	"testing"
)

func always(int) bool { return true }

func TestHistoryUndoRedo(t *testing.T) {
	var h History[int]
	h.Push(1)
	h.Push(2)

	if v, ok := h.Undo(always); !ok || v != 2 {
		t.Fatalf("Undo() = %d, %v; want 2, true", v, ok)
	}
	if v, ok := h.Redo(); !ok || v != 2 {
		t.Fatalf("Redo() = %d, %v; want 2, true", v, ok)
	}
	if !slices.Equal(h.Passed(), []int{1, 2}) {
		t.Fatalf("Passed() = %v, want [1 2]", h.Passed())
	}
}

func TestHistoryUndoStopsWhenRejected(t *testing.T) {
	var h History[int]
	h.Push(1)
	h.Push(2)

	if _, ok := h.Undo(func(v int) bool { return v != 2 }); ok {
		t.Fatal("Undo() went past a rejected entry")
	}
	if !slices.Equal(h.Passed(), []int{1, 2}) {
		t.Fatalf("Passed() = %v, want [1 2]", h.Passed())
	}
}

func TestHistoryPushDumpsFuture(t *testing.T) {
	var h History[int]
	h.Push(1)
	h.Undo(always)
	h.Push(3)

	if _, ok := h.Redo(); ok {
		t.Fatal("Redo() after Push, want nothing to redo")
	}
	if !slices.Equal(h.Passed(), []int{3}) {
		t.Fatalf("Passed() = %v, want [3]", h.Passed())
	}
}

func TestHistoryEmpty(t *testing.T) {
	var h History[int]

	if _, ok := h.Undo(always); ok {
		t.Fatal("Undo() on an empty history reported true")
	}
	if _, ok := h.Redo(); ok {
		t.Fatal("Redo() on an empty history reported true")
	}
}
