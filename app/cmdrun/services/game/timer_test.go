package game

import "testing"

func TestTimerStartsPausedAndResumes(t *testing.T) {
	timer := Time()

	if history := timer.GetHistory(); len(history) != 0 {
		t.Fatalf("new timer history length = %d, want 0", len(history))
	}
	if timer.GetLastTimeStop() != 0 {
		t.Fatal("new timer has a stop time")
	}

	timer.Start()
	if len(timer.GetHistory()) != 1 {
		t.Fatal("started timer should have one running record")
	}

	timer.Pause()
	history := timer.GetHistory()
	if len(history) != 1 {
		t.Fatalf("paused timer history length = %d, want 1", len(history))
	}
	if history[0].Start == 0 || history[0].Stop < history[0].Start {
		t.Fatalf("invalid paused record: %+v", history[0])
	}
	if timer.GetLastTimeStop() != history[0].Stop {
		t.Fatalf("last stop = %d, want %d", timer.GetLastTimeStop(), history[0].Stop)
	}

	timer.Resume()
	history = timer.GetHistory()
	if len(history) != 2 {
		t.Fatalf("resumed timer history length = %d, want 2", len(history))
	}
	if history[1].Start == 0 || history[1].Stop < history[1].Start {
		t.Fatalf("invalid resumed record: %+v", history[1])
	}
}

func TestTimerIgnoresInvalidStateTransitions(t *testing.T) {
	timer := Time()

	timer.Pause()
	timer.Resume()
	if timer.init {
		t.Fatal("pause or resume initialized a new timer")
	}

	timer.Start()
	timer.Start()
	if len(timer.timeStart) != 1 {
		t.Fatalf("repeated start created %d starts, want 1", len(timer.timeStart))
	}

	timer.Pause()
	timer.Pause()
	if len(timer.timeRun) != 1 {
		t.Fatalf("repeated pause created %d runs, want 1", len(timer.timeRun))
	}

	timer.Resume()
	timer.Resume()
	if len(timer.timeStart) != 2 {
		t.Fatalf("repeated resume created %d starts, want 2", len(timer.timeStart))
	}
}
