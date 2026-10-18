package session

import "testing"

func TestStopwatchStartsPausedAndResumes(t *testing.T) {
	stopwatch := NewStopwatch()

	if history := stopwatch.GetHistory(); len(history) != 0 {
		t.Fatalf("new stopwatch history length = %d, want 0", len(history))
	}
	if stopwatch.GetLastTimeStop() != 0 {
		t.Fatal("new stopwatch has a stop time")
	}

	stopwatch.Start()
	if len(stopwatch.GetHistory()) != 1 {
		t.Fatal("started stopwatch should have one running record")
	}

	stopwatch.Pause()
	history := stopwatch.GetHistory()
	if len(history) != 1 {
		t.Fatalf("paused stopwatch history length = %d, want 1", len(history))
	}
	if history[0].Start == 0 || history[0].Stop < history[0].Start {
		t.Fatalf("invalid paused record: %+v", history[0])
	}
	if stopwatch.GetLastTimeStop() != history[0].Stop {
		t.Fatalf("last stop = %d, want %d", stopwatch.GetLastTimeStop(), history[0].Stop)
	}

	stopwatch.Resume()
	history = stopwatch.GetHistory()
	if len(history) != 2 {
		t.Fatalf("resumed stopwatch history length = %d, want 2", len(history))
	}
	if history[1].Start == 0 || history[1].Stop < history[1].Start {
		t.Fatalf("invalid resumed record: %+v", history[1])
	}
}

func TestStopwatchIgnoresInvalidStateTransitions(t *testing.T) {
	stopwatch := NewStopwatch()

	stopwatch.Pause()
	stopwatch.Resume()
	if stopwatch.init {
		t.Fatal("pause or resume initialized a new stopwatch")
	}

	stopwatch.Start()
	stopwatch.Start()
	if len(stopwatch.timeStart) != 1 {
		t.Fatalf("repeated start created %d starts, want 1", len(stopwatch.timeStart))
	}

	stopwatch.Pause()
	stopwatch.Pause()
	if len(stopwatch.timeRun) != 1 {
		t.Fatalf("repeated pause created %d runs, want 1", len(stopwatch.timeRun))
	}

	stopwatch.Resume()
	stopwatch.Resume()
	if len(stopwatch.timeStart) != 2 {
		t.Fatalf("repeated resume created %d starts, want 2", len(stopwatch.timeStart))
	}
}
