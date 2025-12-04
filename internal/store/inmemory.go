package store

import (
	"andi-custodian/internal/chain"
	"context"
	"sync"
)

type InMemoryStore struct {
	mu        sync.RWMutex
	transfers map[string]*TransferResult
	nonces    map[string]uint64
	utxos     map[string][]chain.UTXO
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		transfers: make(map[string]*TransferResult),
		nonces:    make(map[string]uint64),
		utxos:     make(map[string][]chain.UTXO),
	}
}

func (s *InMemoryStore) GetTransferResult(ctx context.Context, id string) (*TransferResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if res, ok := s.transfers[id]; ok {
		return res, nil
	}
	return nil, nil // or return a sentinel error
}

func (s *InMemoryStore) SaveTransferResult(ctx context.Context, id string, result *TransferResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transfers[id] = result
	return nil
}

func (s *InMemoryStore) GetNonce(ctx context.Context, address string) (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nonces[address], nil
}

func (s *InMemoryStore) SetNonce(ctx context.Context, address string, nonce uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nonces[address] = nonce
	return nil
}

func (s *InMemoryStore) GetUTXOs(ctx context.Context, address string) ([]chain.UTXO, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	utxos, ok := s.utxos[address]
	if !ok {
		return nil, nil
	}
	// Return a copy to prevent mutation
	result := make([]chain.UTXO, len(utxos))
	copy(result, utxos)
	return result, nil
}

func (s *InMemoryStore) SaveUTXOs(ctx context.Context, address string, utxos []chain.UTXO) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Store a copy
	s.utxos[address] = make([]chain.UTXO, len(utxos))
	copy(s.utxos[address], utxos)
	return nil
}
