// Package wallet provides secure key management and derivation for multi-chain custody.
package wallet

import (
	"github.com/tyler-smith/go-bip39"
)

// GenerateMnemonic creates a new 12-word BIP-39 mnemonic for HD wallet recovery.
// It is the root of all derived keys and must be stored securely by the user.
func GenerateMnemonic() (string, []byte, error) {
	// 128 bits = 12-word mnemonic
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", nil, err
	}
	seed := bip39.NewSeed(mnemonic, "") // no passphrase
	return mnemonic, seed, nil
}

// NewWallet creates a new wallet from a BIP-39 mnemonic.
// The seed is derived deterministically and used for HD key derivation.
func NewWallet(mnemonic string) (*Wallet, error) {
	seed := bip39.NewSeed(mnemonic, "")
	return &Wallet{seed: seed}, nil
}

// Wallet holds the seed and provides methods to derive addresses and sign.
type Wallet struct {
	seed []byte
}
