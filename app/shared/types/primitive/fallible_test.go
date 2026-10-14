package primitive

import (
	"errors"
	"testing"
)

func TestFallibleErrors(t *testing.T) {
	var f Fallible

	if f.HasErr() || f.GetErr() != nil {
		t.Fatal("new Fallible unexpectedly has an error")
	}

	first := errors.New("first")
	second := errors.New("second")
	f.SetErr(nil)
	f.SetErr(first)
	f.SetErr(second)

	if !f.HasErr() {
		t.Fatal("Fallible has no error after SetErr")
	}
	if !errors.Is(f.GetErr(), first) || !errors.Is(f.GetErr(), second) {
		t.Fatal("GetErr does not contain all set errors")
	}

	f.ClearErr()
	if f.HasErr() || f.GetErr() != nil {
		t.Fatal("Fallible still has errors after ClearErr")
	}
}
