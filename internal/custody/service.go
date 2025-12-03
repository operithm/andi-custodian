package custody

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"andi-custodian/internal/wallet"
)

// Chain represents supported blockchains
type Chain string

const (
	EthereumSepolia Chain = "ethereum-sepolia"
	BitcoinTestnet  Chain = "bitcoin-testnet"
)

// TransferRequest defines a cross-chain custody transfer
type TransferRequest struct {
	ID        string // idempotency key
	Chain     Chain
	ToAddress string
	Amount    string // e.g., "1.5 ETH", "0.05 BTC" (parsed later)
	Payload   []byte // raw tx hash to sign
}

// TransferResult captures outcome
type TransferResult struct {
	TxID      string
	Status    string // "pending", "confirmed", "failed"
	BlockHash string
	Timestamp time.Time
}

// NonceManager tracks Ethereum nonces safely
type NonceManager struct {
	mu     sync.Mutex
	nonces map[string]uint64 // address -> next nonce
}

func NewNonceManager() *NonceManager {
	return &NonceManager{
		nonces: make(map[string]uint64),
	}
}

func (nm *NonceManager) Next(address string) uint64 {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	n := nm.nonces[address]
	nm.nonces[address] = n + 1
	return n
}

// UTXO represents an unspent output (simplified)
type UTXO struct {
	TxID  string
	VOut  uint32
	Value int64 // satoshis
}

// UTXOSelector picks UTXOs to cover a target amount
type UTXOSelector struct{}

func (us *UTXOSelector) Select(utxos []UTXO, target int64) ([]UTXO, int64, error) {
	// Greedy selection (largest first)
	total := int64(0)
	var selected []UTXO
	for _, u := range utxos {
		if total >= target {
			break
		}
		selected = append(selected, u)
		total += u.Value
	}
	if total < target {
		return nil, 0, errors.New("insufficient funds")
	}
	return selected, total - target, nil // change = total - target
}

// TransferService orchestrates multi-chain transfers
type TransferService struct {
	signer       wallet.Signer
	nonceManager *NonceManager
	utxoSelector *UTXOSelector
	idempotency  sync.Map // requestID -> *TransferResult
}

func NewTransferService(signer wallet.Signer) *TransferService {
	return &TransferService{
		signer:       signer,
		nonceManager: NewNonceManager(),
		utxoSelector: &UTXOSelector{},
	}
}

// Transfer initiates a custody transfer with idempotency
func (s *TransferService) Transfer(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	// 1. Idempotency check
	if existing, ok := s.idempotency.Load(req.ID); ok {
		return existing.(*TransferResult), nil
	}

	// 2. Sign the payload (simulated MPC)
	sig, err := s.signer.Sign(ctx, wallet.SignRequest{
		Chain:   wallet.Chain(req.Chain),
		Payload: req.Payload,
	})
	if err != nil {
		return nil, fmt.Errorf("signing failed: %w", err)
	}

	// 3. Broadcast would happen here (simulated)
	txID := fmt.Sprintf("mock-tx-%x", sig[:8])

	result := &TransferResult{
		TxID:      txID,
		Status:    "pending",
		Timestamp: time.Now(),
	}

	// 4. Store result for idempotency
	s.idempotency.Store(req.ID, result)

	return result, nil
}
