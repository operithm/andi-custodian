// nft_test.go
package wallet

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
	"math/big"
	"testing"
)

func TestSimulatedMPCSigner_SignNFTTransfer_ERC721(t *testing.T) {
	mnemonic := "slab lonely fish push bomb festival open oval empower federal slot hotel"
	seed := bip39.NewSeed(mnemonic, "")
	signer := NewSimulatedMPCSigner(seed)

	req := NFTTransferRequest{
		Chain:    EthereumSepolia,
		From:     "0x8E76C1897e55d208b2b5f45cDb43FD7d403a9a31",
		To:       "0x742d35Cc6634C0532925a3b844Bc9dbd8b5E8a18",
		Contract: "0x5d3a536E4D6DbD6114cc1Ead35777bAB948E3643", // example NFT
		TokenID:  big.NewInt(12345),
		Standard: ERC721,
	}

	sig, err := signer.SignNFTTransfer(context.Background(), req)
	assert.NoError(t, err)
	assert.Len(t, sig, 65) // Ethereum signature
}

func TestSimulatedMPCSigner_SignNFTTransfer_ERC1155(t *testing.T) {
	mnemonic := "slab lonely fish push bomb festival open oval empower federal slot hotel"
	seed := bip39.NewSeed(mnemonic, "")
	signer := NewSimulatedMPCSigner(seed)

	req := NFTTransferRequest{
		Chain:    EthereumSepolia,
		From:     "0x8E76C1897e55d208b2b5f45cDb43FD7d403a9a31",
		To:       "0x742d35Cc6634C0532925a3b844Bc9dbd8b5E8a18",
		Contract: "0x495f947276749Ce646f68AC8c248420045cb7b5e", // example ERC-1155
		TokenID:  big.NewInt(1000000),
		Value:    big.NewInt(1), // transfer 1 unit
		Standard: ERC1155,
	}

	sig, err := signer.SignNFTTransfer(context.Background(), req)
	assert.NoError(t, err)
	assert.Len(t, sig, 65)
}

func TestSimulatedMTPSigner_SignNFTTransfer_Ordinals(t *testing.T) {
	mnemonic := "slab lonely fish push bomb festival open oval empower federal slot hotel"
	seed := bip39.NewSeed(mnemonic, "")
	signer := NewSimulatedMPCSigner(seed)

	req := NFTTransferRequest{
		Chain:    BitcoinTestnet,
		From:     "tb1q4d750u3s88c6mt8732j2q6gsn23rwwey25xxnm",
		To:       "tb1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq",
		Contract: "inscription-id-123", // simulated
		TokenID:  big.NewInt(1),
		Value:    big.NewInt(1), // 1 sat
		Standard: ORDINALS,
	}

	sig, err := signer.SignNFTTransfer(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, sig) // Bitcoin DER signature (70-72 bytes)
}
