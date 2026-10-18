package cmdserverunsession

import (
	"context"
	pluminterface "github.com/cyril-fl/plum/interface"
	"os"
	"os/signal"

	cmdservetui "solopg/app/cmdserve/tui"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango/repository"
	"solopg/app/shared/services/process/generatesteplist"

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

type ProcessTemplate[T pluminterface.Stringable] struct {
	GetRepository func() repository.Watchable[T]
	GetFlags      func() *pflag.FlagSet
}

func Process[T pluminterface.Stringable](
	params ProcessTemplate[T],
) *run[T] {
	return &run[T]{
		getRepository: params.GetRepository,
		getFlags:      params.GetFlags,
	}
}

func (p *run[T]) Run() {
	p.makeRepository()
	p.makeTuiView()
	p.runTuiView()
}

// Getters & Setters
func (p *run[T]) GetResult() {
	logs.SilentWarning("error:not_implemented", map[string]any{
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

	/*
		FIXME MEDIUM devrais être une vue et nom l'equilavent du moteur qui tourne dans CMDRUN
		Deplaver le moteur aiderais a ameliore et fixer ici
		Ne devrai pas rester comme ça ca ca marche mais ma melanger les responsabilté
	*/
	model := cmdservetui.New(cmdservetui.StreamModelTemplate[T]{
		Context:    ctx,
		Steps:      p.cache.view.GetSteps(),
		Repository: p.repository,
		Flags:      p.getFlags(),
	})

	p.cache.view.Run(&model, tea.WithContext(ctx))
	if model.HasErr() {
		p.SetErr(model.GetErr())
	}

	if p.cache.view.HasErr() {
		p.SetErr(p.cache.view.GetErr())
	}
}

// Helpers
