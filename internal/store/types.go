// internal/store/types.go
package store

import "time"

// TransferResult represents the outcome of a custody transfer.
type TransferResult struct {
	TxID      string    `json:"tx_id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
