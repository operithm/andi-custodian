// types_test.go
package chain

import (
	"math/big"
	"testing"
)

func TestTxRequest_Validation(t *testing.T) {
	req := &TxRequest{
		Chain: BitcoinTestnet,
		From:  "tb1q4d750u3s88c6mt8732j2q6gsn23rwwey25xxnm",
		To:    "tb1q4d750u3s88c6mt8732j2q6gsn23rwwey25xxnm",
		Value: big.NewInt(1000000), // 10 BTC in satoshis
		ID:    "req-123",
	}
	if req.Value.Int64() != 1000000 {
		t.Errorf("Expected Value=1000000, got %d", req.Value.Int64())
	}
}

func TestChain_EnumValues(t *testing.T) {
	if BitcoinTestnet != "bitcoin-testnet" {
		t.Error("BitcoinTestnet enum value changed")
	}
	if EthereumSepolia != "ethereum-sepolia" {
		t.Error("EthereumSepolia enum value changed")
	}
}
