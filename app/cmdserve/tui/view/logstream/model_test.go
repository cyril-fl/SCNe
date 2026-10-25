package logstream

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/pflag"

	"scne/app/shared/services/logs"
)

type fakeWatchable struct {
	stream chan logs.Log
}

func (f *fakeWatchable) Watch(context.Context) (<-chan logs.Log, error) {
	return f.stream, nil
}

func TestViewRendersStreamedLogs(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("filter", "", "")
	first := logs.Log{Type: logs.INFO, Author: "test", Message: "first"}
	repo := &fakeWatchable{stream: make(chan logs.Log, 1)}
	repo.stream <- first

	m := NewModel(UiParams[logs.Log]{
		Context:    context.Background(),
		Repository: repo,
		Flags:      flags,
	})

	_, wait := m.Update(m.Init()())
	m.Update(wait())

	if content := m.View().Content; !strings.Contains(content, first.String()) {
		t.Fatalf("view does not contain streamed log: %q", content)
	}
}
