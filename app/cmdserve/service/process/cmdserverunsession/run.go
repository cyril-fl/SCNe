package cmdserverunsession

import (
	"context"
	pluminterface "github.com/cyril-fl/plum/interface"
	"os"
	"os/signal"

	cmdservetui "scne/app/cmdserve/tui"
	"scne/app/cmdserve/tui/view/logstream"
	"scne/app/shared/components/models"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/mango/repository"
	"scne/app/shared/services/process/generatesteplist"

	tea "charm.land/bubbletea/v2"
	plumprimitive "github.com/cyril-fl/plum/primitive"
	"github.com/spf13/pflag"
)

type run[T pluminterface.Stringable] struct {
	plumprimitive.Process
	cache[T]

	getRepository func() repository.Watchable[T]
	getFlags      func() *pflag.FlagSet
}

type cache[T pluminterface.Stringable] struct {
	view       *generatesteplist.Generator[cmdservetui.Context]
	repository repository.Watchable[T]
}

// ProcessTemplate provides the watched repository and the command flags.
type ProcessTemplate[T pluminterface.Stringable] struct {
	GetRepository func() repository.Watchable[T]
	GetFlags      func() *pflag.FlagSet
}

// Process returns a process streaming the repository in a TUI.
func Process[T pluminterface.Stringable](
	params ProcessTemplate[T],
) *run[T] {
	return &run[T]{
		getRepository: params.GetRepository,
		getFlags:      params.GetFlags,
	}
}

// Run runs the log stream TUI until interrupted.
func (p *run[T]) Run() {
	p.makeRepository()
	p.makeTuiView()
	p.runTuiView()
}

// Getters & Setters

// GetResult is not implemented; it only logs a warning.
func (p *run[T]) GetResult() {
	logs.SilentWarning("error.not_implemented", map[string]any{
		"Function": "GetResult",
		"Subject":  "session",
	})
}

// Methods
func (p *run[T]) makeRepository() {
	p.repository = p.getRepository()
}

func (p *run[T]) makeTuiView() {
	p.cache.view = generatesteplist.Process[cmdservetui.Context]()
}

func (p *run[T]) runTuiView() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	p.cache.view.Add(models.Step[cmdservetui.Context]{
		Submodel: logstream.NewModel(logstream.UiParams[T]{
			Context:    ctx,
			Repository: p.repository,
			Flags:      p.getFlags(),
		}),
	})

	model := cmdservetui.NewModel(p.cache.view.GetSteps())

	p.cache.view.Run(&model, tea.WithContext(ctx))

	if p.cache.view.HasErr() {
		p.SetErr(p.cache.view.GetErr())
	}
}

// Helpers
