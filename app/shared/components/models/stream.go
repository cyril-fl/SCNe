package models

import (
	"context"
	"strings"

	"scne/app/shared/messages"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/mango/repository"

	pluminterface "github.com/cyril-fl/plum/interface"

	tea "charm.land/bubbletea/v2"
)

// StreamModelTemplate configures a StreamModel.
type StreamModelTemplate[T pluminterface.Stringable] struct {
	Context    context.Context
	Repository repository.Watchable[T]
}

// StreamModel collects the events watched from a repository.
type StreamModel[T pluminterface.Stringable] struct {
	ctx        context.Context
	repository repository.Watchable[T]

	events []T
	stream <-chan T
}

// NewStreamModel returns a model watching params.Repository.
func NewStreamModel[T pluminterface.Stringable](params StreamModelTemplate[T]) *StreamModel[T] {
	return &StreamModel[T]{
		ctx:        params.Context,
		repository: params.Repository,
	}
}

// Init starts watching the repository.
func (m *StreamModel[T]) Init() tea.Cmd {
	return startWatch(m.ctx, m.repository)
}

// Update appends each event and waits for the next one.
// It resolves with an error when the watch fails or the stream closes.
func (m *StreamModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case streamReadyMsg[T]:
		m.stream = msg.stream
		return m, waitForEvent(m.stream)
	case streamEventMsg[T]:
		m.events = append(m.events, msg.event)
		return m, waitForEvent(m.stream)

	case streamErrorMsg:
		return m, messages.Resolve(messages.Resolved{Err: msg.err})
	case streamClosedMsg:
		err := logs.CeaseError("error.unexpected:mongo:close")
		return m, messages.Resolve(messages.Resolved{Err: err})
	}

	return m, nil
}

// View renders one event per line.
func (m *StreamModel[T]) View() tea.View {
	lines := make([]string, 0, len(m.events))
	for _, event := range m.events {
		lines = append(lines, event.String())
	}

	return tea.NewView(strings.Join(lines, "\n"))
}

// Getters & Setters

// GetEvents returns the received events.
func (m *StreamModel[T]) GetEvents() []T {
	return m.events
}

// Helpers
func startWatch[T any](ctx context.Context, repo repository.Watchable[T]) tea.Cmd {
	return func() tea.Msg {
		stream, err := repo.Watch(ctx)
		if err != nil {
			return streamErrorMsg{err: err}
		}

		return streamReadyMsg[T]{stream: stream}
	}
}

func waitForEvent[T any](stream <-chan T) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-stream
		if !ok {
			return streamClosedMsg{}
		}

		return streamEventMsg[T]{event: event}
	}
}

type streamReadyMsg[T any] struct {
	stream <-chan T
}

type streamEventMsg[T any] struct {
	event T
}

type streamErrorMsg struct {
	err error
}

type streamClosedMsg struct{}
