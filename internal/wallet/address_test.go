// address_test.go
package wallet

import (
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

// Use a fixed mnemonic for deterministic addresses
const testMnemonic = "slab lonely fish push bomb festival open oval empower federal slot hotel"

func TestDeriveAddress_BitcoinTestnet(t *testing.T) {
	wallet, err := NewWallet(testMnemonic)
	if err != nil {
		t.Fatal(err)
	}
	addr, err := wallet.DeriveAddress(BitcoinTestnet)
	if err != nil {
		t.Fatal(err)
	}
	s := addr.(string)
	// Should be Bech32 testnet (starts with tb1)
	if len(s) == 0 || s[:3] != "tb1" {
		t.Errorf("Invalid Bitcoin testnet address: %s", s)
	}
}

func TestDeriveAddress_EthereumSepolia(t *testing.T) {
	wallet, err := NewWallet(testMnemonic)
	if err != nil {
		t.Fatal(err)
	}
	addr, err := wallet.DeriveAddress(EthereumSepolia)
	if err != nil {
		t.Fatal(err)
	}
	e := addr.(common.Address).String() // or use common.Address
	if len(e) != 42 || e[:2] != "0x" {
		t.Errorf("Invalid Ethereum address: %s", e)
	}
}
