// threshold.go
package wallet

import (
	"crypto/sha256"
	"errors"
)

// ThresholdPolicy defines how many shares are required to sign.
type ThresholdPolicy struct {
	Threshold int // e.g., 2
	Total     int // e.g., 3
}

// Share represents a simulated MPC key share (in real MPC, this would be a partial key).
type Share struct {
	ID  int
	Key []byte // simulated; real MPC would use secret sharing
}

// ThresholdKey simulates an MPC-managed key with policy.
type ThresholdKey struct {
	Policy  ThresholdPolicy
	Shares  []Share
	RootKey []byte // for simulation only — NEVER in production
}

// NewThresholdKeyFromSeed creates a simulated threshold key from a BIP-39 seed.
// In real MPC, this would be done via distributed key generation (DKG).
func NewThresholdKeyFromSeed(seed []byte, policy ThresholdPolicy) (*ThresholdKey, error) {
	if policy.Threshold < 1 || policy.Threshold > policy.Total {
		return nil, errors.New("invalid threshold policy")
	}

	// For simulation: use first 32 bytes as root key
	rootKey := seed[:32]
	if len(rootKey) != 32 {
		return nil, errors.New("seed too short")
	}

	// Simulate shares (in real MPC, these would be generated via DKG)
	shares := make([]Share, policy.Total)
	for i := 0; i < policy.Total; i++ {
		// Deterministic "share" for demo (not cryptographically secure!)
		h := sha256.Sum256(append(rootKey, byte(i)))
		shares[i] = Share{ID: i + 1, Key: h[:32]}
	}

	return &ThresholdKey{
		Policy:  policy,
		Shares:  shares,
		RootKey: rootKey,
	}, nil
}
