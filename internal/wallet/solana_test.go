// wallet/solana_test.go
package wallet

import (
	"context"
	"github.com/tyler-smith/go-bip39"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeriveSolanaAddress(t *testing.T) {
	mnemonic := "slab lonely fish push bomb festival open oval empower federal slot hotel"
	seed := bip39.NewSeed(mnemonic, "")
	addr, err := DeriveSolanaAddress(seed)
	assert.NoError(t, err)
	assert.Len(t, addr, 44) // Solana addresses are 44 chars base58
}

func TestSimulatedMPCSigner_SignSolana(t *testing.T) {
	mnemonic := "slab lonely fish push bomb festival open oval empower federal slot hotel"
	seed := bip39.NewSeed(mnemonic, "")
	signer := NewSimulatedMPCSigner(seed)

	msg := []byte("hello solana")
	sig, err := signer.SignSolana(context.Background(), msg)
	assert.NoError(t, err)
	assert.Len(t, sig, 64) // Ed25519 signatures are 64 bytes
}
