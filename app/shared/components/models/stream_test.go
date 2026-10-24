package models

import (
	"context"
	"errors"
	"testing"

	"solopg/app/shared/messages"
)

type testEvent string

func (e testEvent) String() string { return string(e) }

type fakeWatchable[T any] struct {
	stream <-chan T
	err    error
	calls  int
	ctx    context.Context
}

func (f *fakeWatchable[T]) Watch(ctx context.Context) (<-chan T, error) {
	f.calls++
	f.ctx = ctx
	return f.stream, f.err
}

func newTestStream(ctx context.Context, repo *fakeWatchable[testEvent]) *StreamModel[testEvent] {
	return NewStreamModel(StreamModelTemplate[testEvent]{
		Context:    ctx,
		Repository: repo,
	})
}

func TestStreamModelReceivesEventsWithoutOpeningMultipleStreams(t *testing.T) {
	stream := make(chan testEvent, 1)
	stream <- "42"

	repo := &fakeWatchable[testEvent]{stream: stream}
	m := newTestStream(context.Background(), repo)

	ready := m.Init()()
	_, wait := m.Update(ready)
	event := wait()
	_, next := m.Update(event)

	if got := m.GetEvents(); len(got) != 1 || got[0] != "42" {
		t.Fatalf("unexpected events: %v", got)
	}
	if repo.calls != 1 {
		t.Fatalf("Watch called %d times, want 1", repo.calls)
	}
	if next == nil {
		t.Fatal("expected a command waiting for the next event")
	}
}

func TestStreamModelResolvesErrorWhenWatchFails(t *testing.T) {
	want := errors.New("watch failed")
	m := newTestStream(context.Background(), &fakeWatchable[testEvent]{err: want})

	_, cmd := m.Update(m.Init()())

	resolved, ok := cmd().(messages.Resolved)
	if !ok || !errors.Is(resolved.Err, want) {
		t.Fatalf("got %v, want Resolved with %v", resolved, want)
	}
}

func TestStreamModelResolvesErrorWhenStreamCloses(t *testing.T) {
	stream := make(chan testEvent)
	close(stream)
	m := newTestStream(context.Background(), &fakeWatchable[testEvent]{stream: stream})

	_, wait := m.Update(m.Init()())
	_, cmd := m.Update(wait())

	resolved, ok := cmd().(messages.Resolved)
	if !ok || resolved.Err == nil {
		t.Fatal("expected a Resolved error when the stream closes")
	}
}

func TestStreamModelViewDoesNotOpenTheRepositoryStream(t *testing.T) {
	repo := &fakeWatchable[testEvent]{}
	m := newTestStream(context.Background(), repo)

	_ = m.View()

	if repo.calls != 0 {
		t.Fatalf("View called Watch %d times", repo.calls)
	}
}

func TestStreamModelViewRendersEvents(t *testing.T) {
	m := newTestStream(context.Background(), &fakeWatchable[testEvent]{})
	m.events = []testEvent{"first", "second"}

	if got := m.View().Content; got != "first\nsecond" {
		t.Fatalf("view = %q", got)
	}
}

func TestStreamModelUsesTheProvidedContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := &fakeWatchable[testEvent]{stream: make(chan testEvent)}
	m := newTestStream(ctx, repo)
	_ = m.Init()()

	if repo.ctx != ctx {
		t.Fatal("Watch did not receive the model context")
	}
}
