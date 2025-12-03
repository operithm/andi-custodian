// ethereum_test.go
package chain

import (
	"errors"
	"math/big"
	"testing"
)

func TestEthereumBuilder_BuildTx_Success(t *testing.T) {
	builder := &EthereumBuilder{}
	req := &TxRequest{
		Chain: EthereumSepolia,
		From:  "0x8E76C1897e55d208b2b5f45cDb43FD7d403a9a31",
		To:    "0x8E76C1897e55d208b2b5f45cDb43FD7d403a9a31",
		Value: big.NewInt(1000000000000000000), // 1 ETH
		ID:    "test-eth",
	}

	opts := BuildOptions{Nonce: 5}
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

func TestEthereumBuilder_BuildTx_InvalidAddress(t *testing.T) {
	builder := &EthereumBuilder{}
	req := &TxRequest{
		Chain: EthereumSepolia,
		From:  "invalid",
		To:    "0x8E76C1897e55d208b2b5f45cDb43FD7d403a9a31",
		Value: big.NewInt(1),
	}

	opts := BuildOptions{Nonce: 0}
	_, err := builder.BuildTx(req, opts)
	if !errors.Is(err, ErrInvalidAddress) {
		t.Errorf("Expected ErrInvalidAddress, got %v", err)
	}
}
