// chain.go
package chain

import "errors"

// Builder abstracts transaction construction across blockchains.
// It returns an *unsigned* transaction for signing by the wallet layer.
type Builder interface {
	BuildTx(req *TxRequest, opts BuildOptions) (*TxResult, error)
}

// NewBuilder creates a chain-specific builder.
func NewBuilder(chainType Chain) (Builder, error) {
	switch chainType {
	case BitcoinTestnet:
		return &BitcoinBuilder{}, nil
	case EthereumSepolia:
		return &EthereumBuilder{}, nil
	default:
		return nil, errors.New("unsupported chain: " + string(chainType))
	}
}
