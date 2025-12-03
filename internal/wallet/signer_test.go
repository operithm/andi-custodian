// signer_test.go
package wallet

import (
	"context"
	"crypto/rand"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

// signer_test.go
func TestSimulatedMPCSigner_Sign_Ethereum(t *testing.T) {
	// Use seed[:32] to derive key — same as signer.go
	seed := bip39.NewSeed(testMnemonic, "")
	privKey, _ := btcec.PrivKeyFromBytes(seed[:32])
	expectedAddr := crypto.PubkeyToAddress(privKey.ToECDSA().PublicKey).Hex()

	signer := NewSimulatedMPCSigner(seed)
	payload := make([]byte, 32)
	rand.Read(payload)

	sig, err := signer.Sign(context.Background(), SignRequest{
		Chain:   EthereumSepolia,
		Payload: payload,
	})
	if err != nil {
		t.Fatalf("ETH sign failed: %v", err)
	}

	verifier := &Verifier{}
	if !verifier.VerifyEthereum(payload, sig, expectedAddr) {
		t.Error("ETH signature verification failed")
	}
}

func TestSimulatedMPCSigner_Sign_Bitcoin(t *testing.T) {
	mnemonic := "slab lonely fish push bomb festival open oval empower federal slot hotel"
	wallet, _ := NewWallet(mnemonic)
	signer := NewSimulatedMPCSigner(wallet.seed)

	payload := make([]byte, 32)
	rand.Read(payload)

	sig, err := signer.Sign(context.Background(), SignRequest{
		Chain:   BitcoinTestnet,
		Payload: payload,
	})
	if err != nil {
		t.Fatalf("BTC sign failed: %v", err)
	}
	if len(sig) < 70 || len(sig) > 72 {
		t.Errorf("BTC signature length = %d, expected 70-72", len(sig))
	}

	// Verify via Verifier
	pubKey, _ := GetPublicKeyFromSeed(wallet.seed, BitcoinTestnet) // helper
	verifier := &Verifier{}
	if !verifier.VerifyBitcoin(payload, sig, pubKey) {
		t.Error("BTC signature verification failed")
	}
}
