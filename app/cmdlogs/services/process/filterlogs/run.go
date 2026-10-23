package filterlogs

import (
	"fmt"
	"solopg/app/shared/services/factory/millfilterlog"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango/repository"

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

type ProcessTemplate struct {
	GetList  func(filter bson.M, opts ...repository.FindOption) ([]logs.Log, error)
	GetFlags func() *pflag.FlagSet
}

func Process(params ProcessTemplate) *run {
	return &run{
		getList:  params.GetList,
		getFlags: params.GetFlags,
	}
}

func (p *run) Run() {
	p.makeList()
	p.makeMill()
	p.displayResult()
}

// Getters & Setters
func (p *run) GetResult() {
	logs.SilentWarning("error:not_implemented", map[string]any{
		"Function": "GetResult",
		"Subject":  "session",
	})
}

// Methods
func (p *run) makeList() {
	list, err := p.getList(nil)
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
		fmt.Println(log.String())
	}
}

// Helpers
