package menuaction

import (
	"os"
	"testing"

	"solopg/app/shared/services/logs"
)

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func TestMain(m *testing.M) {
	logs.Init(testLogRegistrar{})
	os.Exit(m.Run())
}
