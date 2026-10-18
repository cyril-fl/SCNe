package cmdservetui

import (
	"context"
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	pluminterface "github.com/cyril-fl/plum/interface"
	"github.com/spf13/pflag"

	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango/repository"
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

func newTestModel[T pluminterface.Stringable](ctx context.Context, repo repository.Watchable[T]) model[T] {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("filter", "", "")

	return New(StreamModelTemplate[T]{
		Context:    ctx,
		Repository: repo,
		Flags:      flags,
	})
}

func TestModelReceivesEventsWithoutOpeningMultipleStreams(t *testing.T) {
	stream := make(chan testEvent, 1)
	stream <- "42"

	repo := &fakeWatchable[testEvent]{stream: stream}
	m := newTestModel(context.Background(), repo)

	ready := m.Init()()
	_, wait := m.Update(ready)
	event := wait()
	_, next := m.Update(event)

	if got := m.events; len(got) != 1 || got[0] != "42" {
		t.Fatalf("unexpected events: %v", got)
	}
	if repo.calls != 1 {
		t.Fatalf("Watch called %d times, want 1", repo.calls)
	}
	if next == nil {
		t.Fatal("expected a command waiting for the next event")
	}
}

func TestModelReturnsErrorWhenWatchFails(t *testing.T) {
	want := errors.New("watch failed")
	repo := &fakeWatchable[testEvent]{err: want}
	m := newTestModel(context.Background(), repo)

	msg := m.Init()()
	_, quit := m.Update(msg)

	if !errors.Is(m.GetErr(), want) {
		t.Fatalf("model error = %v, want %v", m.GetErr(), want)
	}
	if _, ok := quit().(tea.QuitMsg); !ok {
		t.Fatal("expected tea.Quit after Watch failure")
	}
}

func TestModelReturnsErrorWhenStreamCloses(t *testing.T) {
	stream := make(chan testEvent)
	close(stream)
	repo := &fakeWatchable[testEvent]{stream: stream}
	m := newTestModel(context.Background(), repo)

	ready := m.Init()()
	_, wait := m.Update(ready)
	closed := wait()
	_, quit := m.Update(closed)

	if m.GetErr() == nil {
		t.Fatal("expected an error when the stream closes")
	}
	if _, ok := quit().(tea.QuitMsg); !ok {
		t.Fatal("expected tea.Quit after stream closure")
	}
}

func TestViewDoesNotOpenTheRepositoryStream(t *testing.T) {
	repo := &fakeWatchable[testEvent]{}
	m := newTestModel(context.Background(), repo)

	_ = m.View()

	if repo.calls != 0 {
		t.Fatalf("View called Watch %d times", repo.calls)
	}
}

func TestViewRendersEvents(t *testing.T) {
	repo := &fakeWatchable[logs.Log]{}
	m := newTestModel(context.Background(), repo)
	first := logs.Log{Type: logs.INFO, Author: "test", Message: "first"}
	second := logs.Log{Type: logs.INFO, Author: "test", Message: "second"}
	m.events = []logs.Log{first, second}

	content := m.View().Content
	if !strings.Contains(content, first.String()) || !strings.Contains(content, second.String()) {
		t.Fatalf("view does not contain events: %q", content)
	}
}

func TestModelQuitsOnCtrlQ(t *testing.T) {
	m := newTestModel(context.Background(), &fakeWatchable[testEvent]{})

	_, quit := m.Update(tea.KeyPressMsg(tea.Key{Code: 'q', Mod: tea.ModCtrl}))

	if _, ok := quit().(tea.QuitMsg); !ok {
		t.Fatal("expected tea.Quit after q")
	}
}

func TestModelUsesTheProvidedContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := &fakeWatchable[testEvent]{stream: make(chan testEvent)}
	m := newTestModel(ctx, repo)
	_ = m.Init()()

	if repo.ctx != ctx {
		t.Fatal("Watch did not receive the model context")
	}

	cancel()
	select {
	case <-repo.ctx.Done():
	default:
		t.Fatal("provided context was not cancelled")
	}
}
