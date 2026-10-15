package millfilterlog

import (
	"testing"

	"solopg/app/shared/services/logs"

	"github.com/spf13/pflag"
)

func makeTestFlags(t *testing.T, filter string, tail int) *pflag.FlagSet {
	t.Helper()
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("filter", filter, "")
	flags.Int("tail", tail, "")
	return flags
}

func TestFilterAndTail(t *testing.T) {
	list := []logs.Log{{Type: logs.INFO}, {Type: logs.ERR}, {Type: logs.INFO}}
	m := New(Template{LogList: list, Flags: makeTestFlags(t, "info", 1)})
	m.Filter().Tail()

	if len(m.GetList()) != 1 || m.GetList()[0].Type != logs.INFO {
		t.Fatalf("filtered list = %#v, want one info log", m.GetList())
	}
}

func TestInvalidFilterLeavesListUnchanged(t *testing.T) {
	list := []logs.Log{{Type: logs.INFO}}
	m := New(Template{LogList: list, Flags: makeTestFlags(t, "unknown", 0)})
	m.Filter()
	if len(m.GetList()) != 1 {
		t.Fatal("invalid filter changed the list")
	}
}
