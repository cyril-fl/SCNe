package cmdserverunsession

import (
	"testing"

	"scne/app/shared/services/mango/repository"

	"github.com/spf13/pflag"
)

type testStringable string

func (v testStringable) String() string { return string(v) }

func TestProcessStoresFactories(t *testing.T) {
	getRepository := func() repository.Watchable[testStringable] { return nil }
	getFlags := func() *pflag.FlagSet { return pflag.NewFlagSet("test", pflag.ContinueOnError) }
	p := Process(ProcessTemplate[testStringable]{GetRepository: getRepository, GetFlags: getFlags})

	if p.getRepository == nil || p.getFlags == nil {
		t.Fatal("process did not retain its factories")
	}
}
