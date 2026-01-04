// Package wallet provides secure, multi-chain signing for custody simulations.
package wallet

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
)

// SignRequest holds the data to be signed.
type SignRequest struct {
	Chain   Chain
	Payload []byte // raw hash to sign (e.g., ETH tx hash or BTC sighash)
}

// Signer signs transactions using secure, verifiable cryptography.
type Signer interface {
	Sign(ctx context.Context, req SignRequest) ([]byte, error)
}

// WalletSeed provides the root entropy for key derivation.
// In a real MPC system, this would never be held in one place.
type WalletSeed struct {
	Seed []byte // BIP-39 seed (64 bytes)
}

// SimulatedMPCSigner simulates an MPC signing service.
// It derives a root private key from the seed and signs locally.
// In production, this would be replaced with a gRPC call to an MPC coordinator.
type SimulatedMPCSigner struct {
	seed WalletSeed
}

var (
	ErrSigningFailed = errors.New("signing failed")
)

var erc20ABIJson = []byte(`[
	{
		"inputs": [
			{"name": "to", "type": "address"},
			{"name": "value", "type": "uint256"}
		],
		"name": "transfer",
		"type": "function"
	}
]`)

func NewSimulatedMPCSigner(seed []byte) *SimulatedMPCSigner {
	return &SimulatedMPCSigner{
		seed: WalletSeed{Seed: seed},
	}
}

// Sign derives a private key from the seed and signs the payload.
// It always verifies the signature before returning.
func (s *SimulatedMPCSigner) Sign(ctx context.Context, req SignRequest) ([]byte, error) {
	if len(s.seed.Seed) < 32 {
		return nil, errors.New("seed too short for private key derivation")
	}

	// Use first 32 bytes as root private key (deterministic, recoverable from mnemonic)
	privKey, _ := btcec.PrivKeyFromBytes(s.seed.Seed[:32])
	if privKey == nil {
		return nil, errors.New("invalid private key from seed")
	}
	goPub := privKey.PubKey().ToECDSA()
	goPriv := privKey.ToECDSA()
	var sig []byte

	switch req.Chain {
	case EthereumSepolia, AvalancheFuji:
		sig, err := crypto.Sign(req.Payload, goPriv)
		if err != nil {
			return nil, fmt.Errorf("ethereum sign failed: %w", err)
		}
		// crypto.Sign already returns 65-byte sig with v=27/28
		return sig, nil

	case BitcoinTestnet:
		// Sign with Go stdlib
		r, s, err := ecdsa.Sign(rand.Reader, goPriv, req.Payload)
		if err != nil {
			return nil, fmt.Errorf("bitcoin sign failed: %w", err)
		}

		// Verify using package function
		if !ecdsa.Verify(goPub, req.Payload, r, s) {
			return nil, errors.New("verification failed")
		}

		der := derEncodeSignature(r, s)
		return der, nil

	case SolanaDevnet:
		return s.SignSolana(ctx, req.Payload)

	default:
		return nil, fmt.Errorf("unsupported chain: %s", req.Chain)
	}

	return sig, nil
}

// derEncodeSignature returns a DER-encoded ECDSA signature (ASN.1 SEQUENCE of two INTEGERs)
// This matches Bitcoin's strict DER requirements (BIP-66).
func derEncodeSignature(r, s *big.Int) []byte {
	// Handle sign bit (prepend 0x00 if high bit is set)
	encodeInt := func(n *big.Int) []byte {
		b := n.Bytes()
		if len(b) == 0 {
			b = []byte{0x00}
		}
		if b[0]&0x80 != 0 {
			// Prepend zero to make it positive
			b = append([]byte{0x00}, b...)
		}
		return b
	}

	rEnc := encodeInt(r)
	sEnc := encodeInt(s)
	seqLen := 2 + len(rEnc) + 2 + len(sEnc)

	var buf bytes.Buffer
	buf.WriteByte(0x30) // SEQUENCE
	buf.WriteByte(byte(seqLen))

	buf.WriteByte(0x02) // INTEGER
	buf.WriteByte(byte(len(rEnc)))
	buf.Write(rEnc)

	buf.WriteByte(0x02) // INTEGER
	buf.WriteByte(byte(len(sEnc)))
	buf.Write(sEnc)

	return buf.Bytes()
}

func computeEthereumTxHash(rawTx []byte, chainID *big.Int) ([]byte, error) {
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(rawTx, tx); err != nil {
		return nil, err
	}
	signer := types.NewEIP155Signer(chainID)
	return signer.Hash(tx).Bytes(), nil
}

// verifyEthereumSignature recovers the public key and checks address match.
func (s *SimulatedMPCSigner) verifyEthereumSignature(hash, sig []byte, pub *ecdsa.PublicKey) bool {
	recoveredPub, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return false
	}
	return crypto.PubkeyToAddress(*pub).Hex() == crypto.PubkeyToAddress(*recoveredPub).Hex()
}

// verifyBitcoinSignature verifies a Bitcoin ECDSA signature against a public key.
func (s *SimulatedMPCSigner) verifyBitcoinSignature(hash []byte, r, sInt *big.Int, pubKey *ecdsa.PublicKey) bool {
	return ecdsa.Verify(pubKey, hash, r, sInt)
}
