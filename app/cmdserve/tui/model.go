package cmdservetui

import (
	"context"
	pluminterface "github.com/cyril-fl/plum/interface"
	plumprimitive "github.com/cyril-fl/plum/primitive"
	plumtypes "github.com/cyril-fl/plum/types"

	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango/repository"
	"solopg/app/shared/types/step"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/pflag"
)

type Context struct{}

type contextStepList = step.List[Context]
type Step = step.Step[Context]

type StreamModelTemplate[T pluminterface.Stringable] struct {
	Context    context.Context
	Steps      *contextStepList
	Repository repository.Watchable[T]
	Flags      *pflag.FlagSet
}

type model[T pluminterface.Stringable] struct {
	plumprimitive.Fallible

	ctx        context.Context
	repository repository.Watchable[T]
	events     []T
	stream     <-chan T
	flags      *pflag.FlagSet

	errorset plumtypes.Set[error]

	size *tea.WindowSizeMsg
}

func New[T pluminterface.Stringable](params StreamModelTemplate[T]) model[T] {
	return model[T]{
		ctx:        params.Context,
		repository: params.Repository,
		flags:      params.Flags,
		errorset:   plumtypes.New[error](),
	}
}

func (m *model[T]) Init() tea.Cmd {
	return startWatch(*m)
}

func (m *model[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case streamReadyMsg[T]:
		m.stream = msg.stream
		m.assertFlag()
		return m, waitForEvent(m.stream)
	case streamEventMsg[T]:
		m.events = append(m.events, msg.event)
		return m, waitForEvent(m.stream)

	case streamErrorMsg:
		m.SetErr(msg.err)
		return m, tea.Quit
	case streamClosedMsg:
		m.SetErr(logs.CeaseError("error.unexpected:mongo:close"))
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.size = &msg
		return m, nil
	case tea.KeyMsg:
		switch plumtypes.Key(msg.String()) {
		case plumtypes.SHORT_CTRL_C, plumtypes.SHORT_CTRL_Q:
			return m, tea.Quit
		}
	}

	return m, nil
}

// Helpers
func startWatch[T pluminterface.Stringable](m model[T]) tea.Cmd {
	return func() tea.Msg {
		stream, err := m.repository.Watch(m.ctx)
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
