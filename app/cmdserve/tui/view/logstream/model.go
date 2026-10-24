package logstream

import (
	"context"

	"solopg/app/shared/components/models"
	"solopg/app/shared/services/mango/repository"

	pluminterface "github.com/cyril-fl/plum/interface"
	"github.com/spf13/pflag"

	tea "charm.land/bubbletea/v2"
)

type model[T pluminterface.Stringable] struct {
	stream *models.StreamModel[T]
	flags  *pflag.FlagSet
}

type UiParams[T pluminterface.Stringable] struct {
	Context    context.Context
	Repository repository.Watchable[T]
	Flags      *pflag.FlagSet
}

func NewModel[T pluminterface.Stringable](params UiParams[T]) *model[T] {
	stream := models.NewStreamModel(models.StreamModelTemplate[T]{
		Context:    params.Context,
		Repository: params.Repository,
	})

	return &model[T]{
		stream: stream,
		flags:  params.Flags,
	}
}

func (m *model[T]) Init() tea.Cmd {
	m.assertFilterFlag()
	return m.stream.Init()
}

func (m *model[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	_, cmd = m.stream.Update(msg)
	return m, cmd
}
