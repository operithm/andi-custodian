// nonce_manager.go
package custody

import (
	"sync"
)

// NonceManager safely assigns nonces for Ethereum addresses.
// It supports concurrent requests and tolerates external transactions.
type NonceManager struct {
	mu     sync.RWMutex
	nonces map[string]uint64 // address -> next expected nonce
}

// NewNonceManager creates a new nonce manager.
func NewNonceManager() *NonceManager {
	return &NonceManager{
		nonces: make(map[string]uint64),
	}
}

// GetNext returns the next nonce to use for an address.
// It is safe for concurrent use.
func (nm *NonceManager) GetNext(address string) uint64 {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	n := nm.nonces[address]
	nm.nonces[address] = n + 1
	return n
}

// Reset allows external systems to sync the nonce (e.g., after observing a confirmed tx).
func (nm *NonceManager) Reset(address string, nonce uint64) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	if nonce > nm.nonces[address] {
		nm.nonces[address] = nonce
	}
}
