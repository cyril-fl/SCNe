package id

import (
	"testing"

	"github.com/google/uuid"
)

func TestIDBase64(t *testing.T) {
	value, err := uuid.Parse("e15a0887-fbd2-4800-9b83-a1dbe734a552")
	if err != nil {
		t.Fatalf("uuid.Parse() error = %v", err)
	}

	if got, want := ID(value).Base64(), "4VoIh/vSSACbg6Hb5zSlUg=="; got != want {
		t.Fatalf("Base64() = %q, want %q", got, want)
	}
}
