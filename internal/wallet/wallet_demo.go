package wallet

import (
	"context"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tyler-smith/go-bip39"
)

func RunCustodyDemo(rpcURL string) error {
	// 1. BIP-39
	fmt.Println("\n1. Generating BIP-39 mnemonic (12 words)...")
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed(mnemonic, "")
	fmt.Printf("Mnemonic: %s\n", mnemonic)

	// 2. Bitcoin
	fmt.Println("\n2. Deriving Bitcoin Testnet (Bech32) address...")
	btcAddr, err := deriveHDAddress(seed, "bitcoin")
	if err != nil {
		return fmt.Errorf("BTC: %w", err)
	}
	fmt.Printf("Bitcoin Address (Testnet): %s\n", btcAddr)

	// 3. Ethereum
	fmt.Println("\n3. Deriving Ethereum (Sepolia) address...")
	ethAddrAny, err := deriveHDAddress(seed, "ethereum")
	if err != nil {
		return fmt.Errorf("ETH: %w", err)
	}
	ethAddr := ethAddrAny.(common.Address)
	fmt.Printf("Ethereum Address: %s\n", ethAddr.Hex())

	// 4. UTXO
	fmt.Println("\n4. Simulating UTXO selection (target: 1.15 BTC)...")
	utxos := []struct {
		ID    string
		Value int64
	}{
		{"tx_a", 600_000_000},
	}
	target := int64(115_000_000)
	selected := greedyUTXOSelection(utxos, target)
	if selected != nil {
		total := selected[0].Value
		change := total - target
		fmt.Printf("  → %s (%.6f BTC)\n", selected[0].ID, float64(total)/1e8)
		fmt.Printf("Total selected: %.6f BTC | Change: %.6f BTC\n", float64(total)/1e8, float64(change)/1e8)
	}

	// 5. Nonce
	fmt.Println("\n5. Fetching Ethereum nonce on Sepolia...")
	nonce, err := getEthereumNonce(rpcURL, ethAddr)
	if err != nil {
		log.Printf("⚠️ Nonce fetch failed: %v", err)
	} else {
		fmt.Printf("Next nonce for %s: %d\n", ethAddr.Hex(), nonce)
	}
	return nil
}

// ... (include your working deriveHDAddress, greedyUTXOSelection, getEthereumNonce functions here)
// deriveHDAddress derives a wallet address for the specified chain from a BIP-39 seed.
// Supported chains: "bitcoin", "ethereum"
// Returns:
//   - string (Bech32 testnet address) for "bitcoin"
//   - common.Address for "ethereum"
func deriveHDAddress(seed []byte, chain string) (interface{}, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("new master key: %w", err)
	}

	switch chain {
	case "bitcoin":
		path := []uint32{
			hdkeychain.HardenedKeyStart + 84,
			hdkeychain.HardenedKeyStart + 1,
			hdkeychain.HardenedKeyStart + 0,
			0, 0,
		}
		key := masterKey
		for _, idx := range path {
			key, err = key.Derive(idx)
			if err != nil {
				return nil, fmt.Errorf("BTC derive %d: %w", idx, err)
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

	case "ethereum":
		// Derive using BIP-44: m/44'/60'/0'/0/0
		path := []uint32{
			hdkeychain.HardenedKeyStart + 44,
			hdkeychain.HardenedKeyStart + 60,
			hdkeychain.HardenedKeyStart + 0,
			0,
			0,
		}

		key := masterKey
		for _, idx := range path {
			var err error
			key, err = key.Derive(idx)
			if err != nil {
				return common.Address{}, fmt.Errorf("failed to derive Ethereum key at index %d: %w", idx, err)
			}
		}

		// ✅ CORRECT: Use ECPrivKey() → returns *btcec.PrivateKey
		privKey, err := key.ECPrivKey()
		if err != nil {
			return common.Address{}, fmt.Errorf("failed to extract Ethereum private key: %w", err)
		}

		// Convert to Go stdlib ECDSA key
		goPrivKey := privKey.ToECDSA()
		ethAddr := crypto.PubkeyToAddress(goPrivKey.PublicKey)
		return ethAddr, nil

	default:
		return nil, fmt.Errorf("unsupported chain: %s", chain)
	}
}

// greedyUTXOSelection selects UTXOs in descending order until target is met
func greedyUTXOSelection(utxos []struct {
	ID    string
	Value int64
}, target int64) []struct {
	ID    string
	Value int64
} {
	// Sort descending by value
	for i := 0; i < len(utxos); i++ {
		for j := i + 1; j < len(utxos); j++ {
			if utxos[i].Value < utxos[j].Value {
				utxos[i], utxos[j] = utxos[j], utxos[i]
			}
		}
	}

	var selected []struct {
		ID    string
		Value int64
	}
	var total int64
	for _, u := range utxos {
		if total >= target {
			break
		}
		selected = append(selected, u)
		total += u.Value
	}
	if total < target {
		return nil
	}
	return selected
}

// getEthereumNonce fetches the pending transaction count for an address
func getEthereumNonce(rpcURL string, addr common.Address) (uint64, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return 0, fmt.Errorf("dial RPC: %w", err)
	}
	defer client.Close()

	ctx := context.Background()
	nonce, err := client.PendingNonceAt(ctx, addr)
	if err != nil {
		return 0, fmt.Errorf("get nonce: %w", err)
	}
	return nonce, nil
}
