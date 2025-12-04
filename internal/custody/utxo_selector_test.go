// utxo_selector_test.go
package custody

import (
	"testing"

	"andi-custodian/internal/chain"
)

func TestGreedySelector_Select_Success(t *testing.T) {
	selector := &GreedySelector{}
	utxos := []chain.UTXO{
		{Value: 600_000_000}, // 6 BTC
		{Value: 500_000_000}, // 5 BTC
		{Value: 100_000_000}, // 1 BTC
	}
	target := int64(1_050_000_000) // 10.5 BTC

	selected, change, err := selector.Select(utxos, target)
	if err != nil {
		t.Fatalf("Select failed: %v", err)
	}

	if len(selected) != 2 {
		t.Errorf("Expected 2 UTXOs, got %d", len(selected))
	}
	if change != 50_000_000 {
		t.Errorf("Expected change=50_000_000, got %d", change)
	}
}

func TestGreedySelector_Select_InsufficientFunds(t *testing.T) {
	selector := &GreedySelector{}
	utxos := []chain.UTXO{
		{Value: 100_000_000},
	}
	_, _, err := selector.Select(utxos, 200_000_000)
	if err != chain.ErrInsufficientFunds {
		t.Errorf("Expected ErrInsufficientFunds, got %v", err)
	}
}
