// inmemory_test.go
package store

import (
	"context"
	"sync"
	"testing"
	"time"

	"andi-custodian/internal/chain"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryStore_TransferResult(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	id := "test-transfer-1"

	// Save
	result := &TransferResult{
		TxID:      "tx-123",
		Status:    "pending",
		Timestamp: time.Now(),
	}
	err := store.SaveTransferResult(ctx, id, result)
	assert.NoError(t, err)

	// Retrieve
	retrieved, err := store.GetTransferResult(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, result.TxID, retrieved.TxID)
	assert.Equal(t, result.Status, retrieved.Status)
}

func TestInMemoryStore_Nonce(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	addr := "0x123"

	// Initially 0
	nonce, err := store.GetNonce(ctx, addr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), nonce)

	// Set to 5
	err = store.SetNonce(ctx, addr, 5)
	assert.NoError(t, err)

	// Retrieve updated
	nonce, err = store.GetNonce(ctx, addr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(5), nonce)
}

func TestInMemoryStore_UTXOs(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	addr := "tb1q..."

	utxos := []chain.UTXO{
		{TxID: "tx1", VOut: 0, Value: 1000000},
		{TxID: "tx2", VOut: 1, Value: 2000000},
	}

	// Save
	err := store.SaveUTXOs(ctx, addr, utxos)
	assert.NoError(t, err)

	// Retrieve
	retrieved, err := store.GetUTXOs(ctx, addr)
	assert.NoError(t, err)
	assert.Len(t, retrieved, 2)
	assert.Equal(t, utxos[0].TxID, retrieved[0].TxID)
	assert.Equal(t, utxos[1].Value, retrieved[1].Value)

	// Test copy isolation (mutating retrieved doesn't affect store)
	retrieved[0].Value = 999
	retrieved2, _ := store.GetUTXOs(ctx, addr)
	assert.Equal(t, int64(1000000), retrieved2[0].Value) // unchanged
}

func TestInMemoryStore_Concurrent(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	addr := "0x123"
	const workers = 10

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(nonce uint64) {
			defer wg.Done()
			store.SetNonce(ctx, addr, nonce)
			// No assertion here — just ensure no race crash
		}(uint64(i))
	}
	wg.Wait()

	// Final value is unpredictable in concurrent write, but should not panic
	_, err := store.GetNonce(ctx, addr)
	assert.NoError(t, err)
}
