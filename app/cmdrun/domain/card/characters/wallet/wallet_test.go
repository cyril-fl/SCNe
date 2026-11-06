package wallet

import "testing"

func TestWalletSpendGivesChange(t *testing.T) {
	w := NewWallet(1, 0, 0)

	if err := w.Spend(30); err != nil {
		t.Fatalf("Spend(30) error = %v", err)
	}
	if w != NewWallet(0, 7, 0) {
		t.Fatalf("wallet = %+v, want 0 gold 7 silver 0 copper", w)
	}
}

func TestWalletSpendRefusesWhenTooLow(t *testing.T) {
	w := NewWallet(0, 1, 5)

	if err := w.Spend(16); err == nil {
		t.Fatal("Spend(16) with 15 copper, want an error")
	}
	if w != NewWallet(0, 1, 5) {
		t.Fatalf("wallet = %+v, want unchanged", w)
	}
}

func TestWalletEarnNormalizes(t *testing.T) {
	w := NewWallet(0, 9, 9)
	w.Earn(1)

	if w != NewWallet(1, 0, 0) || w.Total() != 100 {
		t.Fatalf("wallet = %+v, want 1 gold", w)
	}
}
