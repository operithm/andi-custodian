// ethereum.go
package chain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// EthereumBuilder constructs unsigned Ethereum transactions.
type EthereumBuilder struct{}

func (e *EthereumBuilder) BuildTx(req *TxRequest, opts BuildOptions) (*TxResult, error) {
	if req.Chain != EthereumSepolia {
		return nil, errors.New("EthereumBuilder: invalid chain")
	}

	if !common.IsHexAddress(req.From) {
		return nil, fmt.Errorf("invalid from address: %w", ErrInvalidAddress)
	}
	if !common.IsHexAddress(req.To) {
		return nil, fmt.Errorf("invalid to address: %w", ErrInvalidAddress)
	}

	// Create legacy transaction (Sepolia supports EIP-155)
	tx := types.NewTransaction(
		opts.Nonce,
		common.HexToAddress(req.To),
		req.Value,
		21000,                     // gas limit
		big.NewInt(2_000_000_000), // 2 gwei gas price (simulation)
		nil,                       // no data for simple transfer
	)

	// Encode as RLP (unsigned)
	rawTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	// Estimate fee
	gasLimit := big.NewInt(21000)
	gasPrice := big.NewInt(2_000_000_000)
	fee := new(big.Int).Mul(gasLimit, gasPrice)

	return &TxResult{
		RawTx:        rawTx,
		EstimatedFee: fee.Int64(),
	}, nil
}
