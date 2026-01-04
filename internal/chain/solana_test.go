// chain/solana_test.go
package chain

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolanaBuilder_BuildTx_Success(t *testing.T) {
	builder := &SolanaBuilder{}
	req := &TxRequest{
		Chain: SolanaDevnet,
		From:  "7UX2Kk87yue12Gc5cW4Hb6j6g1qX1JvQjG5Y8z9W5tq", // 44 chars
		To:    "9fXw3Kk87yue12Gc5cW4Hb6j6g1qX1JvQjG5Y8z9W8kR",
		Value: big.NewInt(1_000_000_000), // 1 SOL
		ID:    "sol-tx-1",
	}

	result, err := builder.BuildTx(req, BuildOptions{})
	assert.NoError(t, err)
	assert.NotEmpty(t, result.RawTx)
	assert.Equal(t, int64(5000), result.EstimatedFee)
}

func TestSolanaBuilder_BuildTx_InvalidAddress(t *testing.T) {
	builder := &SolanaBuilder{}
	req := &TxRequest{
		Chain: SolanaDevnet,
		From:  "short",
		To:    "9fXw3Kk87yue12Gc5cW4Hb6j6g1qX1JvQjG5Y8z9W8kR",
		Value: big.NewInt(1_000_000_000),
	}

	_, err := builder.BuildTx(req, BuildOptions{})
	assert.Error(t, err)
}
