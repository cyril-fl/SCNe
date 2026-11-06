package primitive

// History keeps entries with undo and redo.
type History[T any] struct {
	passed []T
	future []T // last is the next redo
}

// Getters & Setters

// Passed returns the entries that are not undone.
func (h *History[T]) Passed() []T {
	return h.passed
}

// Methods

// Push adds v and drops the undone entries, which v makes unreachable.
func (h *History[T]) Push(v T) {
	h.DumpFuture()
	h.passed = append(h.passed, v)
}

// Undo moves the last entry to the future when undoable accepts it; otherwise
// it stops and reports false.
func (h *History[T]) Undo(undoable func(T) bool) (T, bool) {
	var zero T
	last := len(h.passed) - 1
	if last < 0 || !undoable(h.passed[last]) {
		return zero, false
	}

	v := h.passed[last]
	h.passed = h.passed[:last]
	h.future = append(h.future, v)
	return v, true
}

// Redo moves the last undone entry back, keeping the rest of the future.
func (h *History[T]) Redo() (T, bool) {
	var zero T
	last := len(h.future) - 1
	if last < 0 {
		return zero, false
	}

	v := h.future[last]
	h.future = h.future[:last]
	h.passed = append(h.passed, v)
	return v, true
}

// DumpFuture drops the undone entries.
func (h *History[T]) DumpFuture() {
	h.future = nil
}
