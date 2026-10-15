package initloggerepository

import (
	"testing"

	"solopg/app/shared/services/mango"
)

func TestProcessStoresDatabase(t *testing.T) {
	db := mango.New()
	p := Process(db)

	if p.db != db {
		t.Fatal("process did not retain its database dependency")
	}
}
