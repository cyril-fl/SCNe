package filterlogs

import (
	"errors"
	"testing"

	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango/repository"

	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/bson"
)

func makeTestFlags(filter string, tail int) func() *pflag.FlagSet {
	return func() *pflag.FlagSet {
		flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
		flags.String("filter", filter, "")
		flags.Int("tail", tail, "")
		return flags
	}
}

func TestRunFiltersAndTailsLoadedList(t *testing.T) {
	list := []logs.Log{{Type: logs.INFO}, {Type: logs.ERR}, {Type: logs.INFO}, {Type: logs.INFO}}
	p := Process(ProcessTemplate{
		GetList:  func(filter bson.M, opts ...repository.FindOption) ([]logs.Log, error) { return list, nil },
		GetFlags: makeTestFlags("info", 2),
	})

	p.Run()

	if p.HasErr() {
		t.Fatalf("Run() error = %v", p.GetErr())
	}
	got := p.mill.GetList()
	if len(got) != 2 || got[0].Type != logs.INFO || got[1].Type != logs.INFO {
		t.Fatalf("result = %#v, want two info logs", got)
	}
}

func TestRunStopsOnListError(t *testing.T) {
	want := errors.New("load failed")
	p := Process(ProcessTemplate{
		GetList: func(filter bson.M, opts ...repository.FindOption) ([]logs.Log, error) { return nil, want },
		GetFlags: func() *pflag.FlagSet {
			t.Fatal("flags requested after list error")
			return nil
		},
	})

	p.Run()

	if !errors.Is(p.GetErr(), want) {
		t.Fatalf("Run() error = %v, want %v", p.GetErr(), want)
	}
	if p.mill != nil {
		t.Fatal("mill built after list error")
	}
}
