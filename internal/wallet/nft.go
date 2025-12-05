// nft.go
package wallet

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// NFTStandard represents supported NFT standards.
type NFTStandard string

const (
	ERC721   NFTStandard = "erc721"
	ERC1155  NFTStandard = "erc1155"
	ORDINALS NFTStandard = "ordinals" // Bitcoin NFTs
)

// NFTTransferRequest defines an NFT transfer.
type NFTTransferRequest struct {
	Chain    Chain
	From     string
	To       string
	Contract string
	TokenID  *big.Int
	Standard NFTStandard
	Value    *big.Int
}

// --- Ethereum: ERC-721 and ERC-1155 ---
var (
	erc721ABIJson  = []byte(`[{"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"transferFrom","type":"function"}]`)
	erc1155ABIJson = []byte(`[{"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"id","type":"uint256"},{"name":"value","type":"uint256"},{"name":"data","type":"bytes"}],"name":"safeTransferFrom","type":"function"}]`)
)

// SignNFTTransfer builds the NFT transfer payload and delegates signing to the embedded Signer.
func (s *SimulatedMPCSigner) SignNFTTransfer(ctx context.Context, req NFTTransferRequest) ([]byte, error) {
	var payload []byte
	var chain Chain

	switch req.Chain {
	case EthereumSepolia:
		var err error
		payload, err = buildEthereumNFTPayload(req)
		if err != nil {
			return nil, err
		}
		chain = EthereumSepolia

	case BitcoinTestnet:
		// For Ordinals: payload = sighash of transaction spending the inscribed UTXO
		// For simulation: use dummy hash (real impl would come from chain/ layer)
		payload = make([]byte, 32) // placeholder
		chain = BitcoinTestnet

	default:
		return nil, fmt.Errorf("unsupported chain for NFT: %s", req.Chain)
	}

	// ✅ DELEGATE TO EXISTING Signer.Sign()
	return s.Sign(ctx, SignRequest{
		Chain:   chain,
		Payload: payload,
	})
}

// buildEthereumNFTPayload constructs the calldata and returns its Keccak256 hash.
func buildEthereumNFTPayload(req NFTTransferRequest) ([]byte, error) {
	from := common.HexToAddress(req.From)
	to := common.HexToAddress(req.To)

	var calldata []byte
	var err error
	switch req.Standard {
	case ERC721:
		abi, _ := abi.JSON(bytes.NewReader(erc721ABIJson))
		calldata, err = abi.Pack("transferFrom", from, to, req.TokenID)
	case ERC1155:
		abi, _ := abi.JSON(bytes.NewReader(erc1155ABIJson))
		calldata, err = abi.Pack("safeTransferFrom", from, to, req.TokenID, req.Value, []byte{})
	default:
		return nil, fmt.Errorf("unsupported Ethereum NFT standard: %s", req.Standard)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to pack NFT call: %w", err)
	}

	return crypto.Keccak256(calldata), nil
}
