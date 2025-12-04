// bitcoin_test.go
package chain

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"math/big"
	"testing"
)

/*
*
const (

	testBtcFrom = "tb1q4d750u3s88c6mt8732j2q6gsn23rwwey25xxnm" // Bech32 -> BIP-173, valid in btcd v0.23.3
	testBtcTo   = "tb1q6r7c6l8v5q7z3a2n4m5k6j7h8g9f0e1d2c3b4a" // Bech32, valid

)
*/
func mustCreateTestnetAddress() string {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		panic(err)
	}
	pubKey := privKey.PubKey()
	addr, err := btcutil.NewAddressWitnessPubKeyHash(
		btcutil.Hash160(pubKey.SerializeCompressed()),
		&chaincfg.TestNet3Params,
	)
	if err != nil {
		panic(err)
	}
	return addr.EncodeAddress()
}

func TestBitcoinBuilder_BuildTx_Success(t *testing.T) {
	fromAddr := mustCreateTestnetAddress()
	toAddr := mustCreateTestnetAddress()
	builder := &BitcoinBuilder{}
	req := &TxRequest{
		Chain: BitcoinTestnet,
		From:  fromAddr,           //valid
		To:    toAddr,             //valid
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
	fromAddr := mustCreateTestnetAddress()
	toAddr := mustCreateTestnetAddress()
	builder := &BitcoinBuilder{}
	req := &TxRequest{
		Chain: BitcoinTestnet,
		From:  fromAddr,            //
		To:    toAddr,              //
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
