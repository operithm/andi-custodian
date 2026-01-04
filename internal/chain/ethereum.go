// ethereum.go
package chain

import (
	_ "andi-custodian/internal/wallet"
	"andi-custodian/pkg/tokens"
	"bytes"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// EthereumBuilder constructs unsigned Ethereum transactions.
type EthereumBuilder struct{}

func (e *EthereumBuilder) BuildTx(req *TxRequest, opts BuildOptions) (*TxResult, error) {
	if req.Chain != EthereumSepolia {
		return nil, errors.New("EthereumBuilder: invalid chain")
	}

	if !common.IsHexAddress(req.From) {
		return nil, fmt.Errorf("invalid from address: %w", ErrInvalidAddress)
	}
	if !common.IsHexAddress(req.To) {
		return nil, fmt.Errorf("invalid to address: %w", ErrInvalidAddress)
	}

	// Create legacy transaction (Sepolia supports EIP-155)
	tx := types.NewTransaction(
		opts.Nonce,
		common.HexToAddress(req.To),
		req.Value,
		21000,                     // gas limit
		big.NewInt(2_000_000_000), // 2 gwei gas price (simulation)
		nil,                       // no data for simple transfer
	)

	// Encode as RLP (unsigned)
	rawTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	// Estimate fee
	gasLimit := big.NewInt(21000)
	gasPrice := big.NewInt(2_000_000_000)
	fee := new(big.Int).Mul(gasLimit, gasPrice)

	return &TxResult{
		RawTx:        rawTx,
		EstimatedFee: fee.Int64(),
	}, nil
}

// BuildTokenTransfer builds an unsigned ERC-20 transfer transaction.
// internal/chain/ethereum.go
func (e *EthereumBuilder) BuildTokenTransfer(req *TokenTransferRequest, nonce uint64) (*TxResult, error) {
	token, ok := tokens.GetTokenBySymbol(string(req.Chain), req.Token)
	if !ok {
		return nil, fmt.Errorf("unsupported token: %s on %s", req.Token, req.Chain)
	}
	if token.IsNative() {
		return nil, errors.New("use BuildTx for native coins")
	}

	amount, err := token.ParseAmount(req.AmountStr)
	if err != nil {
		return nil, err
	}

	// Build ERC-20 transfer calldata
	erc20ABI, _ := abi.JSON(bytes.NewReader(erc20ABIJson))
	calldata, err := erc20ABI.Pack("transfer", common.HexToAddress(req.To), amount)
	if err != nil {
		return nil, err
	}

	// Create transaction to token contract
	tx := types.NewTransaction(
		nonce,
		token.Contract,
		big.NewInt(0),
		65000, // gas limit for ERC-20 transfer
		GetGasPrice(req.Chain),
		calldata,
	)

	rawTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	fee := new(big.Int).Mul(big.NewInt(65000), GetGasPrice(req.Chain))
	return &TxResult{
		RawTx:        rawTx,
		EstimatedFee: fee.Int64(),
	}, nil
}

// GetChainID returns the correct chain ID for EVM chains.
func GetChainID(chainType Chain) *big.Int {
	switch chainType {
	case EthereumSepolia:
		return big.NewInt(11155111) // Sepolia
	case AvalancheFuji:
		return big.NewInt(43113) // Fuji Testnet
	default:
		return big.NewInt(1) // mainnet fallback
	}
}

// GetGasPrice returns a simulated gas price (in wei).
func GetGasPrice(chainType Chain) *big.Int {
	switch chainType {
	case AvalancheFuji:
		return big.NewInt(25_000_000_000) // 25 gwei (Avalanche is cheaper)
	default:
		return big.NewInt(2_000_000_000) // 2 gwei
	}
}

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
