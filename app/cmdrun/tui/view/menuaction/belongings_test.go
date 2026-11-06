package menuaction

import (
	"testing"

	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/card/characters/wallet"
	"scne/app/cmdrun/domain/card/objects"
	"scne/app/cmdrun/domain/gameplay/codex"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/cmdrun/types/id"
)

func newBelongingsSession() *session.GameSession {
	book := codex.New()
	potion := &objects.Object{Category: "potion"}
	potion.Name = "potion"
	book.Objects.Add(potion)
	return session.NewGameSession([16]byte{}, &session.SessionState{
		Codex:  book,
		Player: &characters.Character{Wallet: wallet.NewWallet(0, 1, 0)},
	})
}

func TestPickupAddsQuantity(t *testing.T) {
	gameSession := newBelongingsSession()

	if _, err := pickup(gameSession, gameSession.State.Codex.Objects.Entries[0].ID, 3); err != nil {
		t.Fatalf("pickup error = %v", err)
	}
	if got := gameSession.State.Player.Inventory.Count("potion"); got != 3 {
		t.Fatalf("potions = %d, want 3", got)
	}
	if _, err := pickup(gameSession, id.New(), 1); err == nil {
		t.Fatal("pickup of an unknown entry, want an error")
	}
}

func TestApplyWalletEarnsAndSpends(t *testing.T) {
	gameSession := newBelongingsSession()
	player := gameSession.State.Player

	if _, err := applyWallet(gameSession, map[string]string{"operation": earn, "gold": "1", "copper": "5"}); err != nil {
		t.Fatalf("earn error = %v", err)
	}
	if player.Wallet != wallet.NewWallet(1, 1, 5) {
		t.Fatalf("wallet after earn = %+v, want 1 gp 1 sp 5 cp", player.Wallet)
	}

	if _, err := applyWallet(gameSession, map[string]string{"operation": spend, "silver": "3"}); err != nil {
		t.Fatalf("spend error = %v", err)
	}
	if player.Wallet != wallet.NewWallet(0, 8, 5) {
		t.Fatalf("wallet after spend = %+v, want 0 gp 8 sp 5 cp", player.Wallet)
	}

	if _, err := applyWallet(gameSession, map[string]string{"operation": spend, "gold": "2"}); err == nil || player.Wallet.Total() != 85 {
		t.Fatal("spending more than held must fail and keep the wallet")
	}
}

func TestValidators(t *testing.T) {
	if validatePositive("0") == nil || validatePositive("x") == nil || validatePositive("2") != nil {
		t.Fatal("validatePositive must accept only counts of at least 1")
	}
	if validateCoins("") != nil || validateCoins("0") != nil || validateCoins("-1") == nil {
		t.Fatal("validateCoins must accept empty and counts of at least 0")
	}
}
