package primitive

import "sync"

// Cache stores a lazily loaded value and the error produced while loading it.
// The cache is safe for concurrent reads, loading, updates, and resets.
type Cache[T any] struct {
	mu       sync.RWMutex
	value    T
	loaded   bool
	reported bool
	fallible Fallible
}

// NewCache creates an empty cache.
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{}
}

// Load calls loader once and stores both its value and error. A successful
// loader that returns a zero value still marks the cache as loaded.
func (c *Cache[T]) Load(loader func() (T, error)) (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.loaded {
		return c.value, c.fallible.GetErr()
	}

	value, err := loader()
	c.value = value
	c.loaded = true
	c.fallible.SetErr(err)

	return c.value, err
}

// Value returns the cached value and whether the cache has been loaded.
func (c *Cache[T]) Value() (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.value, c.loaded
}

// Update mutates the cached value under the cache lock. An update makes the
// cache loaded and clears a previous loading error because it replaces the
// failed state with a usable value.
func (c *Cache[T]) Update(update func(*T)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	update(&c.value)
	c.loaded = true
	c.reported = false
	c.fallible.ClearErr()
}

// IsLoaded reports whether the cache has completed loading or has been
// populated through Update.
func (c *Cache[T]) IsLoaded() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.loaded
}

// ReportError invokes report once for the currently stored error. It is
// intended for compatibility wrappers that log an error while exposing a
// simpler value-only API.
func (c *Cache[T]) ReportError(report func(error)) {
	c.mu.Lock()
	if c.reported || !c.fallible.HasErr() {
		c.mu.Unlock()
		return
	}

	err := c.fallible.GetErr()
	c.reported = true
	c.mu.Unlock()

	report(err)
}

// Reset clears the cached value, loading state, and stored errors.
func (c *Cache[T]) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero T
	c.value = zero
	c.loaded = false
	c.reported = false
	c.fallible.ClearErr()
}
