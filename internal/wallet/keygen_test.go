// keygen_test.go
package wallet

import (
	"testing"
)

func TestGenerateMnemonic(t *testing.T) {
	mnemonic, seed, err := GenerateMnemonic()
	if err != nil {
		t.Fatalf("GenerateMnemonic failed: %v", err)
	}
	if len(mnemonic) == 0 {
		t.Error("Mnemonic is empty")
	}
	if len(seed) != 64 {
		t.Errorf("Seed length = %d, want 64", len(seed))
	}
}

func TestNewWallet(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	wallet, err := NewWallet(mnemonic)
	if err != nil {
		t.Fatalf("NewWallet failed: %v", err)
	}
	if wallet.seed == nil {
		t.Error("Seed is nil")
	}
}
