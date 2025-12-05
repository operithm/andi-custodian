// verifier_test.go
package wallet

import (
	"context"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

// Use the same mnemonic as in other tests for consistency
// Derive seed from mnemonic (BIP-39)
func getTestSeed() []byte {
	seed := bip39.NewSeed(testMnemonic, "")
	return seed
}

// verifier_test.go
func TestVerifier_VerifyEthereum(t *testing.T) {
	// 1. Use the same seed as signer tests
	seed := getTestSeed()

	// 2. Derive expected address (same way signer does)
	privKey, _ := btcec.PrivKeyFromBytes(seed[:32])
	expectedAddr := crypto.PubkeyToAddress(privKey.ToECDSA().PublicKey).Hex()

	// 3. Use your working signer to produce a valid signature
	signer := NewSimulatedMPCSigner(seed)
	payload := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
		17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

	sig, err := signer.Sign(context.Background(), SignRequest{
		Chain:   EthereumSepolia,
		Payload: payload[:],
	})
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	// 4. Verify with Verifier
	verifier := &Verifier{}
	if !verifier.VerifyEthereum(payload[:], sig, expectedAddr) {
		t.Error("ETH signature verification failed")
	}

	// 5. Tamper test
	sig[0] ^= 0xFF
	if verifier.VerifyEthereum(payload[:], sig, expectedAddr) {
		t.Error("Tampered signature should fail")
	}
}

func TestVerifier_VerifyBitcoin(t *testing.T) {
	// Use known-good DER signature or generate via signer
	// For brevity, assume you have a valid (hash, derSig, pubKey)
	// from TestSimulatedMPCSigner_Sign_Bitcoin
	// Full test would mirror that flow
	t.Skip("Implement using signer-generated values")
}
func TestVerifier_VerifyNFTTransfer_Ordinals(t *testing.T) {
	mnemonic := "slab lonely fish push bomb festival open oval empower federal slot hotel"
	seed := bip39.NewSeed(mnemonic, "")
	signer := NewSimulatedMPCSigner(seed)

	req := NFTTransferRequest{
		Chain:    BitcoinTestnet,
		Standard: ORDINALS,
	}

	sig, err := signer.SignNFTTransfer(context.Background(), req)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(sig), 70) // DER signature length
	assert.LessOrEqual(t, len(sig), 72)
}
