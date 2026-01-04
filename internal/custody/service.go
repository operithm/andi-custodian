// service.go
package custody

import (
	"andi-custodian/internal/store"
	"andi-custodian/pkg/tokens"
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"andi-custodian/internal/chain"
	"andi-custodian/internal/wallet"
)

// TransferRequest defines a custody transfer.
type TransferRequest struct {
	ID    string
	Chain string
	From  string
	To    string
	Asset string // "ETH", "USTC", "EUTC", "BAYC"
	Value string // "1.0", "1.000000", "12345"
}

// Service orchestrates multi-chain custody operations.
type Service struct {
	signer       wallet.Signer
	store        store.Store   // ← add store dependency
	nonceManager *NonceManager // optional: or delegate to store
	utxoSelector UTXOSelector
	idempotency  sync.Map
}

// NewService creates a new custody service.
func NewService(signer wallet.Signer, store store.Store) *Service {
	return &Service{
		signer:       signer,
		store:        store,
		nonceManager: NewNonceManager(),
		utxoSelector: &GreedySelector{},
	}
}

// Transfer initiates a custody transfer with idempotency.
func (s *Service) Transfer(ctx context.Context, req *TransferRequest) (*store.TransferResult, error) {
	// 1. Idempotency check
	if existing, ok := s.idempotency.Load(req.ID); ok {
		return existing.(*store.TransferResult), nil
	}

	// 2. Parse value (simplified: assume satoshis/wei based on chain)
	// In production, use units package (e.g., 1.5 ETH → 1500000000000000000)
	var valueInt int64 = 1_000_000_000_000_000_000 // 1 ETH or 0.01 BTC for demo

	// 3. Build transaction
	//convert string to chain.Chain
	chainType := chain.Chain(req.Chain)
	builder, err := chain.NewBuilder(chainType)
	if err != nil {
		return nil, err
	}

	_, ok := tokens.GetTokenBySymbol(req.Chain, req.Asset)
	if !ok {
		return nil, fmt.Errorf("unsupported asset: %s on %s", req.Asset, req.Chain)
	}

	var opts chain.BuildOptions
	switch chainType {
	case chain.EthereumSepolia:
		opts.Nonce = s.nonceManager.GetNext(req.From)
	case chain.BitcoinTestnet:
		// In production, fetch UTXOs from indexer or store
		opts.UTXOs = []chain.UTXO{
			{TxID: "mock", VOut: 0, Value: 2000000000}, // 20 BTC
		}
	default:
		return nil, errors.New("unsupported chain")
	}

	txReq := &chain.TxRequest{
		Chain: chainType,
		From:  req.From,
		To:    req.To,
		Value: big.NewInt(valueInt),
		ID:    req.ID,
	}

	tx, err := builder.BuildTx(txReq, opts)
	if err != nil {
		return nil, fmt.Errorf("build tx failed: %w", err)
	}

	// 4. Sign transaction
	// Note: In real system, payload = tx hash (sighash for BTC, keccak256 for ETH)
	sig, err := s.signer.Sign(ctx, wallet.SignRequest{
		Chain:   wallet.Chain(req.Chain),
		Payload: tx.RawTx, // simplified; real system uses tx hash
	})
	if err != nil {
		return nil, fmt.Errorf("signing failed: %w", err)
	}

	// 5. Broadcast would happen here (simulated)
	txID := fmt.Sprintf("mock-tx-%x", sig[:8])

	result := &store.TransferResult{
		TxID:      txID,
		Status:    "pending",
		Timestamp: time.Now(),
	}

	// 6. Store for idempotency
	s.idempotency.Store(req.ID, result)

	// 7. Start monitoring finality (in background)
	go s.monitorFinality(chainType, txID, req.ID)

	return result, nil
}

// monitorFinality simulates finality confirmation.
func (s *Service) monitorFinality(chain chain.Chain, txID, id string) {
	// In production: poll RPC, wait for N confirmations
	time.Sleep(5 * time.Second)

	// Update status (in real system, use a store)
	if existing, ok := s.idempotency.Load(id); ok {
		res := existing.(*store.TransferResult)
		res.Status = "confirmed"
	}
}
