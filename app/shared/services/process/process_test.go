package process

import (
	"errors"
	"testing"
)

type testProcess struct {
	Process
	run func()
}

func (p *testProcess) Run() { p.run() }

func TestHandleProcessStopsAtFirstError(t *testing.T) {
	want := errors.New("failed")
	runs := 0
	first := &testProcess{run: func() { runs++ }}
	second := &testProcess{}
	second.run = func() { runs++; second.SetErr(want) }
	third := &testProcess{run: func() { runs++ }}

	if err := HandleProcess([]Processable{first, second, third}); !errors.Is(err, want) {
		t.Fatalf("HandleProcess error = %v, want %v", err, want)
	}
	if runs != 2 {
		t.Fatalf("ran %d processes, want 2", runs)
	}
}

func TestHandleProcessReturnsNilWhenAllProcessesSucceed(t *testing.T) {
	if err := HandleProcess([]Processable{
		&testProcess{run: func() {}},
		&testProcess{run: func() {}},
	}); err != nil {
		t.Fatalf("HandleProcess error = %v, want nil", err)
	}
}
