// postgres_test.go
package store

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresStore_Exists(t *testing.T) {
	if os.Getenv("TEST_POSTGRES") == "" {
		t.Skip("Skipping PostgreSQL tests (set TEST_POSTGRES=1 to enable)")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=postgres password=postgres dbname=andi_custodian sslmode=disable"
	}

	store, err := NewPostgresStore(connStr)
	require.NoError(t, err)

	// Test nonce
	ctx := context.Background()
	addr := "0x123"
	err = store.SetNonce(ctx, addr, 42)
	assert.NoError(t, err)

	nonce, err := store.GetNonce(ctx, addr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(42), nonce)
}
