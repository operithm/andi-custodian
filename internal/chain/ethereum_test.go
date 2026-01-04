// ethereum_test.go
package chain

import (
	"andi-custodian/pkg/tokens"
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

const (
	EthereumSepoliaFrom = "0x8E76C1897e55d208b2b5f45cDb43FD7d403a9a31"
	EthereumSepoliaTo   = "0x742d35Cc6634C0532925a3b844Bc9dbd8b5E8a18"
)

func TestEthereumBuilder_BuildTx_Success(t *testing.T) {
	builder := &EthereumBuilder{}
	req := &TxRequest{
		Chain: EthereumSepolia,
		From:  EthereumSepoliaFrom,
		To:    EthereumSepoliaFrom,
		Value: big.NewInt(1_000_000_000_000_000_000), // 1 ETH
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
		To:    EthereumSepoliaFrom,
		Value: big.NewInt(1),
	}

	opts := BuildOptions{Nonce: 0}
	_, err := builder.BuildTx(req, opts)
	if !errors.Is(err, ErrInvalidAddress) {
		t.Errorf("Expected ErrInvalidAddress, got %v", err)
	}
}

func TestEthereumBuilder_BuildTokenTransfer(t *testing.T) {
	// Ensure USDC_Sepolia is registered
	token, ok := tokens.GetTokenBySymbol("ethereum-sepolia", "USDC")
	if !ok {
		t.Fatal("USDC not found in token registry")
	}
	expectedContract := token.Contract

	builder := &EthereumBuilder{}

	req := &TokenTransferRequest{
		Chain:     EthereumSepolia,
		From:      EthereumSepoliaFrom,
		To:        EthereumSepoliaTo,
		Token:     "USDC",
		AmountStr: "1.000000", // 1 USDC (6 decimals)
		ID:        "token-tx-1",
	}

	result, err := builder.BuildTokenTransfer(req, 5) // nonce = 5
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.RawTx, "RawTx should not be empty")
	assert.Greater(t, result.EstimatedFee, int64(0), "Fee should be > 0")

	// Optional: Decode and validate transaction
	// (Only if you want deep validation — not required for unit test)
	tx, err := decodeTransaction(result.RawTx)
	assert.NoError(t, err)
	assert.Equal(t, uint64(5), tx.Nonce())
	assert.Equal(t, expectedContract, tx.To(), "Transaction should be sent to token contract")
	assert.Equal(t, big.NewInt(0), tx.Value(), "Token transfers have 0 ETH value")

	// Validate gas limit and price
	assert.Equal(t, uint64(65000), tx.Gas())
	gasPrice := GetGasPrice(EthereumSepolia)
	expectedFee := new(big.Int).Mul(big.NewInt(65000), gasPrice)
	assert.Equal(t, expectedFee.Int64(), result.EstimatedFee)
}

func TestEthereumBuilder_BuildTokenTransfer_InvalidToken(t *testing.T) {
	builder := &EthereumBuilder{}
	req := &TokenTransferRequest{
		Chain:     EthereumSepolia,
		From:      EthereumSepoliaFrom,
		To:        EthereumSepoliaTo,
		Token:     "INVALID_TOKEN",
		AmountStr: "1",
	}

	_, err := builder.BuildTokenTransfer(req, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported token")
}

func TestEthereumBuilder_BuildTx_AvalancheFuji(t *testing.T) {
	builder := &EthereumBuilder{}
	req := &TxRequest{
		Chain: AvalancheFuji,
		From:  EthereumSepoliaFrom,
		To:    EthereumSepoliaTo,
		Value: big.NewInt(1_000_000_000_000_000_000), // 1 AVAX
		ID:    "avax-tx-1",
	}

	result, err := builder.BuildTx(req, BuildOptions{Nonce: 5})
	assert.NoError(t, err)
	assert.NotEmpty(t, result.RawTx)

	// Fee should reflect Avalanche gas price (25 gwei)
	assert.Equal(t, int64(525000), result.EstimatedFee) // 21000 * 25e9
}

// decodeTransaction decodes RLP-encoded transaction bytes.
// Helper for test validation only.
func decodeTransaction(rawTx []byte) (*types.Transaction, error) {
	tx := new(types.Transaction)
	err := rlp.DecodeBytes(rawTx, tx)
	return tx, err
}
