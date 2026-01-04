// types.go
package chain

import (
	"errors"
	"math/big"
)

// Chain represents a supported blockchain.
type Chain string

const (
	BitcoinTestnet  Chain = "bitcoin-testnet"
	EthereumSepolia Chain = "ethereum-sepolia"
	SolanaDevnet    Chain = "solana-devnet"
	AvalancheFuji   Chain = "avalanche-fuji"
)

// TxRequest is a cross-chain transaction request.
type TxRequest struct {
	Chain Chain
	From  string
	To    string
	Value *big.Int
	ID    string // idempotency key
}

// TxResult is the output of transaction building.
type TxResult struct {
	RawTx        []byte // unsigned serialized transaction
	EstimatedFee int64  // in native units (satoshis or wei)
}

// TokenTransferRequest is a cross-chain token transaction request
type TokenTransferRequest struct {
	Chain     Chain
	From      string
	To        string
	Token     string // e.g., "USDC"
	AmountStr string // e.g., "1.000000"
	ID        string // for idempotency
}

// BuildOptions provides chain-specific context (e.g., UTXOs for BTC, nonce for ETH)
type BuildOptions struct {
	UTXOs []UTXO // for Bitcoin
	Nonce uint64 // for Ethereum
}

// UTXO represents an unspent output (Bitcoin only)
type UTXO struct {
	TxID     string
	VOut     uint32
	Value    int64
	PkScript []byte
}

// Errors
var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAddress    = errors.New("invalid address")
)
