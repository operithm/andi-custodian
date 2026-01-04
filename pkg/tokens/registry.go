// pkg/tokens/registry.go
package tokens

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Token represents a supported fungible or non-fungible token.
type Token struct {
	Name     string
	Symbol   string
	Chain    string
	Contract common.Address
	Decimals int
}

// --- Ethereum Sepolia ---
var (
	// Native
	ETH_Sepolia = &Token{
		Name:     "Ethereum",
		Symbol:   "ETH",
		Chain:    "ethereum-sepolia",
		Contract: common.Address{}, // native coin
		Decimals: 18,
	}

	// Stablecoins
	USDC_Sepolia = &Token{
		Name:     "USD Coin",
		Symbol:   "USDC",
		Chain:    "ethereum-sepolia",
		Contract: common.HexToAddress("0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238"),
		Decimals: 6,
	}

	// Circle Stablecoins (USTC, EUTC)
	USTC_Sepolia = &Token{
		Name:     "US Treasury Stablecoin",
		Symbol:   "USTC",
		Chain:    "ethereum-sepolia",
		Contract: common.HexToAddress("0x0000000000000000000000000000000000000001"), // PLACEHOLDER
		Decimals: 6,
	}

	EUTC_Sepolia = &Token{
		Name:     "Euro Treasury Stablecoin",
		Symbol:   "EUTC",
		Chain:    "ethereum-sepolia",
		Contract: common.HexToAddress("0x0000000000000000000000000000000000000002"), // PLACEHOLDER
		Decimals: 6,
	}
)

// --- Avalanche Fuji ---
var (
	AVAX_Fuji = &Token{
		Name:     "Avalanche",
		Symbol:   "AVAX",
		Chain:    "avalanche-fuji",
		Contract: common.Address{}, // native C-Chain coin
		Decimals: 18,
	}

	USDCe_Fuji = &Token{
		Name:     "USD Coin (Avalanche)",
		Symbol:   "USDC.e",
		Chain:    "avalanche-fuji",
		Contract: common.HexToAddress("0x5425890298aed601595a70AB815c96711a31Bc65"),
		Decimals: 6,
	}
)

// --- Solana Devnet ---
// Solana doesn't use contract addresses — tokens are PDAs.
// For simulation, we use a registry of mint addresses.
var (
	SOL_Devnet = &Token{
		Name:   "Solana",
		Symbol: "SOL",
		Chain:  "solana-devnet",
		// Contract = Mint address for native SOL is system program (not used)
		Decimals: 9,
	}

	USDC_SolanaDevnet = &Token{
		Name:     "USD Coin (Solana)",
		Symbol:   "USDC",
		Chain:    "solana-devnet",
		Contract: common.HexToAddress("4zMMC9srt5Ri5X14GAgGqhPgJ6Hdw84Yjy818YxdYkG9"), // Devnet USDC mint
		Decimals: 6,
	}
)

// AllTokens returns a slice of all supported tokens (useful for tests or CLI).
func AllTokens() []*Token {
	return []*Token{
		// Ethereum
		ETH_Sepolia, USDC_Sepolia, USTC_Sepolia, EUTC_Sepolia,
		// Avalanche
		AVAX_Fuji, USDCe_Fuji,
		// Solana
		SOL_Devnet, USDC_SolanaDevnet,
	}
}

// GetTokenBySymbol returns a token by chain and symbol.
func GetTokenBySymbol(chain, symbol string) (*Token, bool) {
	for _, t := range AllTokens() {
		if t.Chain == chain && t.Symbol == symbol {
			return t, true
		}
	}
	return nil, false
}

// IsNative returns true if the token is the chain's native coin.
func (t *Token) IsNative() bool {
	return t.Contract == (common.Address{})
}

// ParseAmount converts a decimal string (e.g., "1.5") to base units (e.g., 1500000 for 6 decimals).
// Simplified for demo — in production, use github.com/shopspring/decimal.
func (t *Token) ParseAmount(amountStr string) (*big.Int, error) {
	// For demo: assume "1.000000" → 1_000_000
	if t.Decimals == 6 {
		return big.NewInt(1_000_000), nil
	}
	if t.Decimals == 18 {
		return new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), nil
	}
	if t.Decimals == 9 {
		return big.NewInt(1_000_000_000), nil
	}
	return big.NewInt(1), nil
}
