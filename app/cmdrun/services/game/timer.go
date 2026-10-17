package game

import (
	"solopg/app/cmdrun/types/timer"
	"time"
)


type Stopwatch struct {
	init      bool
	running   bool
	timeStart []int64 // unix milisecond TODO faire une note
	timeRun   []int64
}

func Time() *Stopwatch {
	return &Stopwatch{
		init:      false,
		running:   false,
		timeStart: []int64{},
		timeRun:   []int64{},
	}
}

// Getters & Setters
func (t *Stopwatch) GetHistory() []timer.Record {
	if !t.init {
		return []timer.Record{}
	}

	history := t.makeHistory()

	if t.running {
		history = append(history, timer.Record{
			Start: getlast(t.timeStart),
			Stop:  getNow(),
		})
	}

	return history
}

func (t *Stopwatch) makeHistory() []timer.Record {
	if !t.init {
		return []timer.Record{}
	}
	maxlength := len(t.timeRun)
	history := make([]timer.Record, maxlength)

	for i, start := range t.timeStart {
		if i >= maxlength {
			break
		}

		history[i] = timer.Record{
			Start: start,
			Stop:  start + t.timeRun[i],
		}
	}
	return history
}

func (t *Stopwatch) GetLastTimeStop() int64 {
	if !t.init {
		return 0
	}
	if t.running {
		return getNow()
	}

	return getlast(t.timeStart) + getlast(t.timeRun)
}

// Methods
func (t *Stopwatch) Start() {
	if t.init {
		return
	}
	t.init = true
	t.running = true
	t.timeStart = append(t.timeStart, getNow())
}

func (t *Stopwatch) Pause() {
	if !t.init {
		return
	}
	if !t.running {
		return
	}
	run := getNow() - getlast(t.timeStart)
	t.timeRun = append(t.timeRun, run)
	t.running = false
}

func (t *Stopwatch) Resume() {
	if !t.init {
		return
	}
	if t.running {
		return
	}
	t.timeStart = append(t.timeStart, getNow())
	t.running = true
}

// Helpers
func getNow() int64 {
	return time.Now().UnixMilli()
}

func getlast[T any](arr []T) T {
	if len(arr) > 0 {
		return arr[len(arr)-1]
	}
	var zero T
	return zero
}

func getRunningTime(history []timer.Record) int64 {
	var total int64
	for _, r := range history {
		total += r.Stop - r.Start
	}
	return total
}
