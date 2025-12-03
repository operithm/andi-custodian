// chain_test.go
package chain

import (
	"testing"
)

func TestNewBuilder_Bitcoin(t *testing.T) {
	builder, err := NewBuilder(BitcoinTestnet)
	if err != nil {
		t.Fatalf("NewBuilder failed: %v", err)
	}
	if _, ok := builder.(*BitcoinBuilder); !ok {
		t.Error("Expected BitcoinBuilder")
	}
}

func TestNewBuilder_Ethereum(t *testing.T) {
	builder, err := NewBuilder(EthereumSepolia)
	if err != nil {
		t.Fatalf("NewBuilder failed: %v", err)
	}
	if _, ok := builder.(*EthereumBuilder); !ok {
		t.Error("Expected EthereumBuilder")
	}
}

func TestNewBuilder_Invalid(t *testing.T) {
	_, err := NewBuilder("invalid-chain")
	if err == nil {
		t.Error("Expected error for invalid chain")
	}
}
