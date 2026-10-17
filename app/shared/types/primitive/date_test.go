package primitive

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	datable.UpdatedAt = updatedAt.Add(-time.Second)
	beforeTouch := datable.GetupdatedAt()

	datable.Touch()

	if !datable.GetcreatedAt().Equal(createdAt) {
		t.Fatal("Touch changed CreatedAt")
	}
	if !datable.GetupdatedAt().After(beforeTouch) {
		t.Fatal("Touch did not advance UpdatedAt")
	}
}

func TestDatableSerializesTimestamps(t *testing.T) {
	datable := Dated()

	encoded, err := bson.Marshal(struct {
		Datable `bson:",inline"`
	}{Datable: datable})
	if err != nil {
		t.Fatalf("bson.Marshal returned an error: %v", err)
	}

	var document bson.M
	if err := bson.Unmarshal(encoded, &document); err != nil {
		t.Fatalf("bson.Unmarshal returned an error: %v", err)
	}

	if _, ok := document["createdAt"]; !ok {
		t.Fatal("BSON document does not contain createdAt")
	}
	if _, ok := document["updatedAt"]; !ok {
		t.Fatal("BSON document does not contain updatedAt")
	}
}

func TestGetAgePreservesDatableState(t *testing.T) {
	datable := Steady()

	age := datable.GetAge()

	if !age.GetcreatedAt().Equal(datable.GetcreatedAt()) {
		t.Fatal("GetAge changed CreatedAt")
	}
	if !age.GetupdatedAt().Equal(datable.GetupdatedAt()) {
		t.Fatal("GetAge changed UpdatedAt")
	}
	if !age.immutable {
		t.Fatal("GetAge changed immutability")
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
