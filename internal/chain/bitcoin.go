// bitcoin.go
package chain

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// BitcoinBuilder constructs unsigned Bitcoin transactions.
type BitcoinBuilder struct{}

func (b *BitcoinBuilder) BuildTx(req *TxRequest, opts BuildOptions) (*TxResult, error) {
	if req.Chain != BitcoinTestnet {
		return nil, errors.New("BitcoinBuilder: invalid chain")
	}

	// Validate addresses
	fromAddr, err := btcutil.DecodeAddress(req.From, &chaincfg.TestNet3Params)
	if err != nil {
		return nil, fmt.Errorf("invalid from address: %w", ErrInvalidAddress)
	}
	toAddr, err := btcutil.DecodeAddress(req.To, &chaincfg.TestNet3Params)
	if err != nil {
		return nil, fmt.Errorf("invalid to address: %w", ErrInvalidAddress)
	}

	// UTXO selection
	target := req.Value.Int64()
	selected, _, change, err := b.selectUTXOs(opts.UTXOs, target)
	if err != nil {
		return nil, err
	}

	// Build transaction
	msgTx := wire.NewMsgTx(wire.TxVersion)

	// Add inputs
	for _, u := range selected {
		txHash, err := chainhash.NewHashFromStr(u.TxID)
		if err != nil {
			return nil, err
		}
		outPoint := wire.NewOutPoint(txHash, u.VOut)
		txIn := wire.NewTxIn(outPoint, nil, nil) // unlocking script will be added during signing
		msgTx.AddTxIn(txIn)
	}

	// Add outputs
	toScript, err := txscript.PayToAddrScript(toAddr)
	if err != nil {
		return nil, err
	}
	msgTx.AddTxOut(wire.NewTxOut(target, toScript))

	if change > 0 {
		changeScript, _ := txscript.PayToAddrScript(fromAddr)
		msgTx.AddTxOut(wire.NewTxOut(change, changeScript))
	}

	// Serialize unsigned tx
	var buf bytes.Buffer
	if err := msgTx.Serialize(&buf); err != nil {
		return nil, err
	}

	// Estimate fee (10 sat/byte)
	fee := int64(len(buf.Bytes())) * 10

	return &TxResult{
		RawTx:        buf.Bytes(),
		EstimatedFee: fee,
	}, nil
}

func (b *BitcoinBuilder) selectUTXOs(utxos []UTXO, target int64) ([]UTXO, int64, int64, error) {
	// Sort descending
	for i := 0; i < len(utxos); i++ {
		for j := i + 1; j < len(utxos); j++ {
			if utxos[i].Value < utxos[j].Value {
				utxos[i], utxos[j] = utxos[j], utxos[i]
			}
		}
	}

	var selected []UTXO
	var total int64
	for _, u := range utxos {
		if total >= target {
			break
		}
		selected = append(selected, u)
		total += u.Value
	}
	if total < target {
		return nil, 0, 0, ErrInsufficientFunds
	}
	return selected, total, total - target, nil
}
