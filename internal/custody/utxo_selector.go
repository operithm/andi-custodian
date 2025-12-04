// utxo_selector.go
package custody

import "andi-custodian/internal/chain"

// UTXOSelector selects UTXOs to fund a transaction.
type UTXOSelector interface {
	Select(utxos []chain.UTXO, target int64) ([]chain.UTXO, int64, error)
}

// GreedySelector implements a simple largest-first selection.
type GreedySelector struct{}

func (gs *GreedySelector) Select(utxos []chain.UTXO, target int64) ([]chain.UTXO, int64, error) {
	// Sort descending by value
	for i := 0; i < len(utxos); i++ {
		for j := i + 1; j < len(utxos); j++ {
			if utxos[i].Value < utxos[j].Value {
				utxos[i], utxos[j] = utxos[j], utxos[i]
			}
		}
	}

	var selected []chain.UTXO
	var total int64
	for _, u := range utxos {
		if total >= target {
			break
		}
		selected = append(selected, u)
		total += u.Value
	}
	if total < target {
		return nil, 0, chain.ErrInsufficientFunds
	}
	return selected, total - target, nil // change = total - target
}
