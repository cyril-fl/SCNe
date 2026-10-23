package mango

import "testing"

func TestDisconnectWithoutClient(t *testing.T) {
	New().Disconnect()
}

func TestInitReturnsConnectionError(t *testing.T) {
	t.Setenv("MONGO_HOST", "localhost")
	t.Setenv("MONGO_PORT", "not-a-port")

	db, err := Init()

	if err == nil {
		t.Fatal("Init() accepted an invalid Mongo URI")
	}
	if db != nil {
		t.Fatalf("Init() = %#v, want nil database on error", db)
	}
}
