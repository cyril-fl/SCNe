package primitive

import (
	"errors"
	"sync"
	"testing"
)

func TestCacheLoadsOnlyOnceIncludingEmptyValues(t *testing.T) {
	cache := NewCache[[]string]()
	calls := 0
	loader := func() ([]string, error) {
		calls++
		return nil, nil
	}

	if value, err := cache.Load(loader); err != nil || value != nil {
		t.Fatalf("first Load() = %#v, %v; want nil, nil", value, err)
	}
	if value, err := cache.Load(loader); err != nil || value != nil {
		t.Fatalf("second Load() = %#v, %v; want nil, nil", value, err)
	}
	if calls != 1 {
		t.Fatalf("loader calls = %d, want 1", calls)
	}
	if !cache.IsLoaded() {
		t.Fatal("empty cache should be marked loaded")
	}
}

func TestCacheStoresAndResetsErrors(t *testing.T) {
	cache := NewCache[string]()
	wantErr := errors.New("load failed")
	calls := 0
	loader := func() (string, error) {
		calls++
		return "", wantErr
	}

	if _, err := cache.Load(loader); !errors.Is(err, wantErr) {
		t.Fatalf("Load() error = %v, want %v", err, wantErr)
	}
	if _, err := cache.Load(func() (string, error) {
		calls++
		return "unexpected", nil
	}); !errors.Is(err, wantErr) {
		t.Fatalf("second Load() error = %v, want %v", err, wantErr)
	}
	if calls != 1 || !cache.fallible.HasErr() || !errors.Is(cache.fallible.GetErr(), wantErr) {
		t.Fatalf("cache error state is inconsistent: calls=%d, hasErr=%t, err=%v", calls, cache.fallible.HasErr(), cache.fallible.GetErr())
	}

	cache.Reset()
	if cache.IsLoaded() || cache.fallible.HasErr() || cache.fallible.GetErr() != nil {
		t.Fatal("Reset() did not clear cache state")
	}
}

func TestCacheUpdateStoresValueAndClearsError(t *testing.T) {
	cache := NewCache[string]()
	wantErr := errors.New("load failed")
	_, _ = cache.Load(func() (string, error) { return "", wantErr })

	cache.Update(func(value *string) { *value = "recovered" })

	if value, loaded := cache.Value(); !loaded || value != "recovered" {
		t.Fatalf("Value() = %q, %t; want recovered, true", value, loaded)
	}
	if cache.fallible.HasErr() || cache.fallible.GetErr() != nil {
		t.Fatal("Update() did not clear the previous error")
	}
}

func TestCacheReportsStoredErrorOnlyOnce(t *testing.T) {
	cache := NewCache[int]()
	wantErr := errors.New("load failed")
	_, _ = cache.Load(func() (int, error) { return 0, wantErr })

	reports := 0
	cache.ReportError(func(err error) {
		reports++
		if !errors.Is(err, wantErr) {
			t.Errorf("reported error = %v, want %v", err, wantErr)
		}
	})
	cache.ReportError(func(error) { reports++ })

	if reports != 1 {
		t.Fatalf("error reports = %d, want 1", reports)
	}
}

func TestCacheLoadsConcurrentlyOnce(t *testing.T) {
	cache := NewCache[int]()
	var mu sync.Mutex
	calls := 0
	loader := func() (int, error) {
		mu.Lock()
		calls++
		mu.Unlock()
		return 42, nil
	}

	var wait sync.WaitGroup
	for i := 0; i < 20; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			if value, err := cache.Load(loader); err != nil || value != 42 {
				t.Errorf("Load() = %d, %v; want 42, nil", value, err)
			}
		}()
	}
	wait.Wait()

	if calls != 1 {
		t.Fatalf("loader calls = %d, want 1", calls)
	}
}
