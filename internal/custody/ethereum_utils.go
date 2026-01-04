// File: internal/custody/ethereum_utils.go
package custody

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

// computeEthereumTxHash computes the EIP-155 transaction hash for signing.
func computeEthereumTxHash(rawTx []byte, chainID *big.Int) ([]byte, error) {
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(rawTx, tx); err != nil {
		return nil, err
	}
	signer := types.NewEIP155Signer(chainID)
	return signer.Hash(tx).Bytes(), nil
}

// decodeTransaction decodes RLP-encoded transaction bytes.
func decodeTransaction(rawTx []byte) (*types.Transaction, error) {
	tx := new(types.Transaction)
	err := rlp.DecodeBytes(rawTx, tx)
	return tx, err
}
