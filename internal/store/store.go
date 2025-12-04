// store/store.go
package store

import (
	"andi-custodian/internal/chain"
	"context"
)

type Store interface {
	// Idempotency
	GetTransferResult(ctx context.Context, id string) (*TransferResult, error)
	SaveTransferResult(ctx context.Context, id string, result *TransferResult) error

	// Ethereum
	GetNonce(ctx context.Context, address string) (uint64, error)
	SetNonce(ctx context.Context, address string, nonce uint64) error

	// Bitcoin
	GetUTXOs(ctx context.Context, address string) ([]chain.UTXO, error)
	SaveUTXOs(ctx context.Context, address string, utxos []chain.UTXO) error
}
