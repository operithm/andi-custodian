// bitcoin_test.go
package chain

import (
	"math/big"
	"testing"
)

func TestBitcoinBuilder_BuildTx_Success(t *testing.T) {
	builder := &BitcoinBuilder{}
	req := &TxRequest{
		Chain: BitcoinTestnet,
		From:  "tb1q4d750u3s88c6mt8732j2q6gsn23rwwey25xxnm",
		To:    "tb1q4d750u3s88c6mt8732j2q6gsn23rwwey25xxnm",
		Value: big.NewInt(500000), // 0.005 BTC
		ID:    "test-btc",
	}

	utxos := []UTXO{
		{TxID: "abc123", VOut: 0, Value: 1000000}, // 0.01 BTC
	}

	opts := BuildOptions{UTXOs: utxos}
	result, err := builder.BuildTx(req, opts)
	if err != nil {
		t.Fatalf("BuildTx failed: %v", err)
	}

	if len(result.RawTx) == 0 {
		t.Error("RawTx is empty")
	}
	if result.EstimatedFee <= 0 {
		t.Error("EstimatedFee must be > 0")
	}
}

func TestBitcoinBuilder_BuildTx_InsufficientFunds(t *testing.T) {
	builder := &BitcoinBuilder{}
	req := &TxRequest{
		Chain: BitcoinTestnet,
		Value: big.NewInt(2000000), // 0.02 BTC
	}

	utxos := []UTXO{
		{TxID: "abc123", VOut: 0, Value: 1000000}, // only 0.01 BTC
	}

	opts := BuildOptions{UTXOs: utxos}
	_, err := builder.BuildTx(req, opts)
	if err != ErrInsufficientFunds {
		t.Errorf("Expected ErrInsufficientFunds, got %v", err)
	}
}

func TestBitcoinBuilder_selectUTXOs(t *testing.T) {
	builder := &BitcoinBuilder{}
	utxos := []UTXO{
		{Value: 600000},
		{Value: 500000},
		{Value: 100000},
	}
	selected, _, change, err := builder.selectUTXOs(utxos, 1050000) // 1.05 BTC
	if err != nil {
		t.Fatalf("selectUTXOs failed: %v", err)
	}
	if len(selected) != 2 {
		t.Errorf("Expected 2 UTXOs, got %d", len(selected))
	}
	if change != 50000 {
		t.Errorf("Expected change=50000, got %d", change)
	}
}
