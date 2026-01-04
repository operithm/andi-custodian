// solana.go
package wallet

import (
	"context"
	"crypto/ed25519"
	"errors"

	"github.com/btcsuite/btcutil/base58"
)

// DeriveSolanaKeypair derives a Solana Ed25519 keypair from the first 32 bytes of the seed.
// Note: Real Solana derivation uses BIP-44 path m/44'/501'/0'/0', but for simulation,
// here use seed[:32] as the secret key seed (sufficient for demo).
func DeriveSolanaKeypair(seed []byte) (ed25519.PrivateKey, error) {
	if len(seed) < 32 {
		return nil, errors.New("seed must be at least 32 bytes")
	}
	secret := seed[:32]
	return ed25519.NewKeyFromSeed(secret), nil
}

// DeriveSolanaAddress returns the Solana public key as a base58-encoded string.
func DeriveSolanaAddress(seed []byte) (string, error) {
	priv, err := DeriveSolanaKeypair(seed)
	if err != nil {
		return "", err
	}
	pub := priv.Public().(ed25519.PublicKey)
	return base58.Encode(pub), nil
}

// SignSolana signs a message using Ed25519.
func (s *SimulatedMPCSigner) SignSolana(ctx context.Context, msg []byte) ([]byte, error) {
	priv, err := DeriveSolanaKeypair(s.seed.Seed)
	if err != nil {
		return nil, err
	}
	return ed25519.Sign(priv, msg), nil
}
