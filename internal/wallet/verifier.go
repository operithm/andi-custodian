// verifier.go
package wallet

import (
	"crypto/ecdsa"
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"
)

type Verifier struct{}

func (v *Verifier) VerifyEthereum(hash, sig []byte, addrHex string) bool {
	if len(sig) != 65 {
		return false
	}
	pub, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return false
	}
	return strings.EqualFold(crypto.PubkeyToAddress(*pub).Hex(), addrHex)
}

func (v *Verifier) VerifyBitcoin(hash, sig []byte, pubKey *btcec.PublicKey) bool {
	// 1. Parse DER signature
	r, s, err := parseSignatureDER(sig)
	if err != nil {
		return false
	}

	// 2. Convert btcec.PublicKey → ecdsa.PublicKey
	goPub := pubKey.ToECDSA()

	// 3. Verify using Go stdlib (you confirmed this works)
	return ecdsa.Verify(goPub, hash, r, s)
}

// parseSignatureDER parses a DER-encoded ECDSA signature into r and s.
// It assumes strict DER (BIP-66).
func parseSignatureDER(sig []byte) (r, s *big.Int, err error) {
	if len(sig) < 8 || sig[0] != 0x30 {
		return nil, nil, errors.New("invalid DER signature: not a SEQUENCE")
	}

	// Skip total length byte
	idx := 2

	// Parse R
	if idx >= len(sig) || sig[idx] != 0x02 {
		return nil, nil, errors.New("invalid DER: R not an INTEGER")
	}
	rLen := int(sig[idx+1])
	idx += 2
	if idx+rLen > len(sig) {
		return nil, nil, errors.New("invalid DER: R length overflow")
	}
	r = new(big.Int).SetBytes(sig[idx : idx+rLen])
	idx += rLen

	// Parse S
	if idx >= len(sig) || sig[idx] != 0x02 {
		return nil, nil, errors.New("invalid DER: S not an INTEGER")
	}
	sLen := int(sig[idx+1])
	idx += 2
	if idx+sLen > len(sig) {
		return nil, nil, errors.New("invalid DER: S length overflow")
	}
	s = new(big.Int).SetBytes(sig[idx : idx+sLen])

	return r, s, nil
}
