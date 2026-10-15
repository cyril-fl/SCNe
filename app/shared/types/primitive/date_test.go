package primitive

import (
	"testing"
	"time"
)

func TestDatedInitializesTimestamps(t *testing.T) {
	datable := Dated()

	if datable.GetcreatedAt().IsZero() || datable.GetupdatedAt().IsZero() {
		t.Fatal("Dated returned zero timestamps")
	}
	if datable.GetcreatedAt().Location() != time.UTC || datable.GetupdatedAt().Location() != time.UTC {
		t.Fatal("Dated timestamps are not in UTC")
	}
	if !datable.GetcreatedAt().Equal(datable.GetupdatedAt()) {
		t.Fatal("Dated timestamps are not initialized to the same instant")
	}
}

func TestDatedTouchUpdatesOnlyUpdatedAt(t *testing.T) {
	datable := Dated()
	createdAt := datable.GetcreatedAt()
	updatedAt := datable.GetupdatedAt()
	datable.updatedAt = updatedAt.Add(-time.Second)
	beforeTouch := datable.GetupdatedAt()

	datable.Touch()

	if !datable.GetcreatedAt().Equal(createdAt) {
		t.Fatal("Touch changed CreatedAt")
	}
	if !datable.GetupdatedAt().After(beforeTouch) {
		t.Fatal("Touch did not advance UpdatedAt")
	}
}

func TestSteadyTouchDoesNotUpdateTimestamp(t *testing.T) {
	datable := Steady()
	updatedAt := datable.GetupdatedAt()

	datable.Touch()

	if !datable.GetupdatedAt().Equal(updatedAt) {
		t.Fatal("Touch changed the timestamp of a steady Datable")
	}
}
