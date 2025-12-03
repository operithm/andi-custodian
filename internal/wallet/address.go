package wallet

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	_ "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Chain represents a supported blockchain network.
type Chain string

const (
	BitcoinTestnet  Chain = "bitcoin-testnet"
	EthereumSepolia Chain = "ethereum-sepolia"
)

// DeriveAddress derives a wallet address for the given chain using standard BIP paths.
func (w *Wallet) DeriveAddress(chain Chain) (interface{}, error) {
	masterKey, err := hdkeychain.NewMaster(w.seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("create master key: %w", err)
	}

	switch chain {
	case BitcoinTestnet:
		path := []uint32{
			hdkeychain.HardenedKeyStart + 84, // BIP-84
			hdkeychain.HardenedKeyStart + 1,  // testnet coin type
			hdkeychain.HardenedKeyStart + 0,
			0, 0,
		}
		key := masterKey
		for _, idx := range path {
			key, err = key.Derive(idx)
			if err != nil {
				return nil, fmt.Errorf("derive BTC key at %d: %w", idx, err)
			}
		}
		pubKey, err := key.ECPubKey()
		if err != nil {
			return nil, err
		}
		addr, err := btcutil.NewAddressWitnessPubKeyHash(
			btcutil.Hash160(pubKey.SerializeCompressed()),
			&chaincfg.TestNet3Params,
		)
		if err != nil {
			return nil, err
		}
		return addr.EncodeAddress(), nil

	case EthereumSepolia:
		path := []uint32{
			hdkeychain.HardenedKeyStart + 44, // BIP-44
			hdkeychain.HardenedKeyStart + 60, // ETH coin type
			hdkeychain.HardenedKeyStart + 0,
			0, 0,
		}
		key := masterKey
		for _, idx := range path {
			key, err = key.Derive(idx)
			if err != nil {
				return nil, fmt.Errorf("derive ETH key at %d: %w", idx, err)
			}
		}
		privKey, err := key.ECPrivKey()
		if err != nil {
			return nil, fmt.Errorf("extract ETH private key: %w", err)
		}
		return crypto.PubkeyToAddress(privKey.ToECDSA().PublicKey), nil

	default:
		return nil, fmt.Errorf("unsupported chain: %s", chain)
	}
}
