package filterlogs

import (
	"scne/app/shared/services/factory/millfilterlog"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/mango/repository"

	"charm.land/lipgloss/v2"
	plumprimitive "github.com/cyril-fl/plum/primitive"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/bson"
)

type run struct {
	plumprimitive.Process
	cache

	getList  func(filter bson.M, opts ...repository.FindOption) ([]logs.Log, error)
	getFlags func() *pflag.FlagSet
}

type cache struct {
	mill *millfilterlog.Mill
	list []logs.Log
}

// ProcessTemplate provides the log source and the command flags.
type ProcessTemplate struct {
	GetList  func(filter bson.M, opts ...repository.FindOption) ([]logs.Log, error)
	GetFlags func() *pflag.FlagSet
}

// Process returns a process printing the logs filtered by the flags.
func Process(params ProcessTemplate) *run {
	return &run{
		getList:  params.GetList,
		getFlags: params.GetFlags,
	}
}

// Run loads the logs, applies --tail then --filter and prints them.
func (p *run) Run() {
	p.makeList()
	p.makeMill()
	p.displayResult()
}

// Getters & Setters

// GetResult is not implemented; it only logs a warning.
func (p *run) GetResult() {
	logs.SilentWarning("error.not_implemented", map[string]any{
		"Function": "GetResult",
		"Subject":  "session",
	})
}

// Methods
func (p *run) makeList() {
	// Without a sort Mongo does not guarantee order, and --tail must keep the newest logs.
	list, err := p.getList(nil, repository.WithSort(bson.D{{Key: "datable.createdAt", Value: 1}}))
	if err != nil {
		p.SetErr(err)
		return
	}

	p.list = list
}

func (p *run) makeMill() {
	if p.HasErr() {
		return
	}

	mill := millfilterlog.New(millfilterlog.Template{
		LogList: p.list,
		Flags:   p.getFlags(),
	})

	if mill.HasErr() {
		p.SetErr(mill.GetErr())
		return
	}

	p.mill = mill
}

func (p *run) displayResult() {
	if p.HasErr() {
		return
	}

	result := p.mill.
		Tail().
		Filter().
		GetList()

	for _, log := range result {
		lipgloss.Println(log.String())
	}
}

// Helpers
