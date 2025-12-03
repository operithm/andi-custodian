// testutil.go
// Package wallet test utilities (used only in *_test.go files).
package wallet

import (
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
)

func GetPublicKeyFromSeed(seed []byte, chain Chain) (*btcec.PublicKey, error) {
	if chain == BitcoinTestnet {
		privKey, _ := btcec.PrivKeyFromBytes(seed[:32])
		if privKey == nil {
			return nil, errors.New("invalid seed for BTC key")
		}
		return privKey.PubKey(), nil // same key as signer!
	}
	// Ethereum can use HD or seed — but for consistency, consider seed too
	return nil, errors.New("unsupported chain")
}
