// solana.go
package chain

import (
	"errors"
)

// SolanaBuilder constructs unsigned Solana transactions (mock for simulation).
type SolanaBuilder struct{}

// BuildTx returns a mock Solana transaction message.
// In production, this would serialize a real Solana Message using solana-go.
func (s *SolanaBuilder) BuildTx(req *TxRequest, opts BuildOptions) (*TxResult, error) {
	if req.Chain != SolanaDevnet {
		return nil, errors.New("SolanaBuilder: invalid chain")
	}

	// Validate addresses (basic length check)
	if len(req.From) != 44 || len(req.To) != 44 {
		return nil, errors.New("invalid Solana address length")
	}

	// For simulation: return a deterministic mock message
	// In production: build real transaction with instructions, recent blockhash, etc.
	mockMsg := []byte("solana-tx-mock-" + req.ID)

	// Estimated fee: 5000 lamports (standard for simple transfer)
	return &TxResult{
		RawTx:        mockMsg,
		EstimatedFee: 5000,
	}, nil
}
