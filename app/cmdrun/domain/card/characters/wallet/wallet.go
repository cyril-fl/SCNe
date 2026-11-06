package wallet

import "scne/app/shared/services/logs"

type Wallet struct {
	Gold   int `yaml:"gold"`
	Silver int `yaml:"silver"`
	Copper int `yaml:"copper"`
}

func NewWallet(gold, silver, copper int) Wallet {
	return Wallet{
		Gold:   gold,
		Silver: silver,
		Copper: copper,
	}
}

// 1 gold = 10 silver = 100 copper.
const (
	copperPerSilver = 10
	copperPerGold   = 100
)

// Total returns the wallet value in copper.
func (w Wallet) Total() int {
	return w.Gold*copperPerGold + w.Silver*copperPerSilver + w.Copper
}

// Earn adds copper and normalizes the coins.
func (w *Wallet) Earn(copper int) {
	w.set(w.Total() + copper)
}

// Spend removes copper, breaking coins to give change; it fails without
// changing the wallet when the total is too low.
func (w *Wallet) Spend(copper int) error {
	if copper > w.Total() {
		return logs.Error("error.invalid:funds", map[string]any{
			"Cost":  copper,
			"Total": w.Total(),
		})
	}
	w.set(w.Total() - copper)
	return nil
}

func (w *Wallet) set(total int) {
	w.Gold = total / copperPerGold
	w.Silver = total % copperPerGold / copperPerSilver
	w.Copper = total % copperPerSilver
}
