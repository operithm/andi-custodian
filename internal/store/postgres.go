// postgres.go
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"andi-custodian/internal/chain"
	_ "github.com/lib/pq"
)

// PostgresStore implements Store using PostgreSQL.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgreSQL store and runs migrations.
func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

// TransferResult methods

func (p *PostgresStore) GetTransferResult(ctx context.Context, id string) (*TransferResult, error) {
	var data []byte
	var timestamp sql.NullTime
	err := p.db.QueryRowContext(ctx,
		"SELECT data, created_at FROM transfers WHERE id = $1",
		id).Scan(&data, &timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var result TransferResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transfer result: %w", err)
	}
	if timestamp.Valid {
		result.Timestamp = timestamp.Time
	}

	return &result, nil
}

func (p *PostgresStore) SaveTransferResult(ctx context.Context, id string, result *TransferResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer result: %w", err)
	}

	_, err = p.db.ExecContext(ctx,
		"INSERT INTO transfers (id, data, created_at) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET data = $2, created_at = $3",
		id, data, result.Timestamp)
	return err
}

// Nonce methods

func (p *PostgresStore) GetNonce(ctx context.Context, address string) (uint64, error) {
	var nonce uint64
	err := p.db.QueryRowContext(ctx,
		"SELECT nonce FROM nonces WHERE address = $1", address).
		Scan(&nonce)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // new address → nonce 0
		}
		return 0, err
	}
	return nonce, nil
}

func (p *PostgresStore) SetNonce(ctx context.Context, address string, nonce uint64) error {
	_, err := p.db.ExecContext(ctx,
		"INSERT INTO nonces (address, nonce) VALUES ($1, $2) ON CONFLICT (address) DO UPDATE SET nonce = $2",
		address, nonce)
	return err
}

// UTXO methods

func (p *PostgresStore) GetUTXOs(ctx context.Context, address string) ([]chain.UTXO, error) {
	rows, err := p.db.QueryContext(ctx,
		"SELECT tx_id, vout, value FROM utxos WHERE address = $1 ORDER BY value DESC",
		address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var utxos []chain.UTXO
	for rows.Next() {
		var u chain.UTXO
		err := rows.Scan(&u.TxID, &u.VOut, &u.Value)
		if err != nil {
			return nil, err
		}
		utxos = append(utxos, u)
	}
	return utxos, rows.Err()
}

func (p *PostgresStore) SaveUTXOs(ctx context.Context, address string, utxos []chain.UTXO) error {
	// Start transaction for atomicity
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing UTXOs for address
	if _, err := tx.ExecContext(ctx, "DELETE FROM utxos WHERE address = $1", address); err != nil {
		return err
	}

	// Insert new UTXOs
	for _, u := range utxos {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO utxos (address, tx_id, vout, value) VALUES ($1, $2, $3, $4)",
			address, u.TxID, u.VOut, u.Value); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Schema
const schema = `
CREATE TABLE IF NOT EXISTS transfers (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS nonces (
    address TEXT PRIMARY KEY,
    nonce BIGINT NOT NULL CHECK (nonce >= 0)
);

CREATE TABLE IF NOT EXISTS utxos (
    address TEXT NOT NULL,
    tx_id TEXT NOT NULL,
    vout INTEGER NOT NULL,
    value BIGINT NOT NULL CHECK (value > 0),
    PRIMARY KEY (address, tx_id, vout)
);

-- Optional: indexes for performance
CREATE INDEX IF NOT EXISTS idx_transfers_id ON transfers(id);
CREATE INDEX IF NOT EXISTS idx_nonces_address ON nonces(address);
CREATE INDEX IF NOT EXISTS idx_utxos_address ON utxos(address);
`
