// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
    "math/big"
    "time"

    "github.com/ponzinomics/go-ethereum/common"
    "github.com/ponzinomics/go-ethereum/rlp"
)

// var (
//	ErrInvalidSig           = errors.New("invalid transaction v, r, s values")
//	ErrUnexpectedProtection = errors.New("transaction type does not supported EIP-155 protected signatures")
//	ErrInvalidTxType        = errors.New("transaction type not valid in this context")
//	ErrTxTypeNotSupported   = errors.New("transaction type not supported")
//	ErrGasFeeCapTooLow      = errors.New("fee cap less than base fee")
//	errShortTypedTx         = errors.New("typed transaction too short")
// )

// Transaction types.
// const (
//	LegacyTxType = iota
//	AccessListTxType
//	DynamicFeeTxType
// )

// TransactionFast is an Ethereum transaction whose only purpose is for analysis in ponzinomics.
type TransactionFast struct {
    time        time.Time
    nonce       uint64
    hash        common.Hash
    data        []byte
    gas         uint64
    gasPrice    *big.Int
    value       *big.Int
    r           *big.Int
    from        string
    to          string
    Type        uint64
    BlockNumber uint64
    chain       uint64
}

// NewTx creates a new transaction.
func NewTxFast(nonce uint64, hash string, gas uint64, gasPrice, value *big.Int, from string, to string) *TransactionFast {
    tx := new(TransactionFast)
    tx.time = time.Now()
    tx.nonce = nonce
    tx.hash = common.HexToHash(hash)
    tx.gas = gas
    tx.gasPrice = new(big.Int).Set(gasPrice)
    tx.value = new(big.Int).Set(value)
    tx.from = from
    tx.to = to
    return tx
}

// EncodeRLP implements rlp.Encoder
// func (tx *TransactionFast) EncodeRLP(w io.Writer) error {
//	if tx.Type == LegacyTxType {
//		return rlp.Encode(w, tx)
//	}
//	// It's an EIP-2718 typed TX envelope.
//	buf := encodeBufferPool.Get().(*bytes.Buffer)
//	defer encodeBufferPool.Put(buf)
//	buf.Reset()
//	if err := tx.encodeTyped(buf); err != nil {
//		return err
//	}
//	return rlp.Encode(w, buf.Bytes())
// }

// encodeTyped writes the canonical encoding of a typed transaction to w.
// func (tx *TransactionFast) encodeTyped(w *bytes.Buffer) error {
//	w.WriteByte(tx.Type)
//	return rlp.Encode(w, tx)
// }

// DecodeRLP implements rlp.Decoder
// func (tx *TransactionFast) UnmarshalBinary(b []byte) error {
//	if len(b) > 0 && b[0] > 0x7f {
//		// It's a legacy transaction.
//		var data LegacyTx
//		err := rlp.DecodeBytes(b, &data)
//		if err != nil {
//			return err
//		}
//		tx.setDecoded(&data, len(b))
//		return nil
//	}
//	// It's an EIP2718 typed transaction envelope.
//	inner, err := tx.decodeTyped(b)
//	if err != nil {
//		return err
//	}
//	tx.setDecoded(inner, len(b))
//	return nil
// }

// decodeTyped decodes a typed transaction from the canonical format.
func (tx *TransactionFast) decodeTyped(b []byte) (TxData, error) {
    if len(b) <= 1 {
        return nil, errShortTypedTx
    }
    switch b[0] {
    case AccessListTxType:
        var inner AccessListTx
        err := rlp.DecodeBytes(b[1:], &inner)
        return &inner, err
    case DynamicFeeTxType:
        var inner DynamicFeeTx
        err := rlp.DecodeBytes(b[1:], &inner)
        return &inner, err
    default:
        return nil, ErrTxTypeNotSupported
    }
}

// ChainId returns the EIP155 chain ID of the transaction. The return value will always be
// non-nil. For legacy transactions which are not replay-protected, the return value is
// zero.
func (tx *TransactionFast) ChainId() uint64 {
    return tx.chain
}

// Data returns the input data of the transaction.
func (tx *TransactionFast) Data() []byte { return tx.data }

// Gas returns the gas limit of the transaction.
func (tx *TransactionFast) Gas() uint64 { return tx.gas }

// GasPrice returns the gas price of the transaction.
func (tx *TransactionFast) GasPrice() *big.Int { return new(big.Int).Set(tx.gasPrice) }

// Value returns the ether amount of the transaction.
func (tx *TransactionFast) Value() *big.Int { return new(big.Int).Set(tx.value) }

// Nonce returns the sender account nonce of the transaction.
func (tx *TransactionFast) Nonce() uint64 { return tx.nonce }

// From returns the recipient address of the transaction.
// For contract-creation transactions, To returns nil.
func (tx *TransactionFast) From() string {
    return tx.from
    // return copyAddressPtr(&tx.to)
}

// To returns the recipient address of the transaction.
// For contract-creation transactions, To returns nil.
func (tx *TransactionFast) To() string {
    return tx.to
    // return copyAddressPtr(&tx.to)
}

func (tx *TransactionFast) Hash() common.Hash {
    return tx.hash
}

// RawSignatureValues returns the V, R, S signature values of the transaction.
// The return values should not be modified by the caller.
func (tx *TransactionFast) RawSignatureValues() (r *big.Int) {
    return new(big.Int).Set(tx.r)
}

// Cost returns gas * gasPrice + value.
func (tx *TransactionFast) Cost() *big.Int {
    total := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(tx.Gas()))
    total.Add(total, tx.Value())
    return total
}

// Transactions implements DerivableList for transactions.
type TransactionsFast []*TransactionFast

// Len returns the length of s.
func (s TransactionsFast) Len() int { return len(s) }

// EncodeIndex encodes the i'th transaction to w. Note that this does not check for errors
// because we assume that *Transaction will only ever contain valid txs that were either
// constructed by decoding or via public API in this package.
// func (s TransactionsFast) EncodeIndex(i int, w *bytes.Buffer) {
//	tx := s[i]
//	if tx.Type() == LegacyTxType {
//		rlp.Encode(w, tx.inner)
//	} else {
//		tx.encodeTyped(w)
//	}
// }

// TxDifference returns a new set which is the difference between a and b.
func TxDifferenceFast(a, b TransactionsFast) TransactionsFast {
    keep := make(TransactionsFast, 0, len(a))

    remove := make(map[common.Hash]struct{})
    for _, tx := range b {
        remove[tx.hash] = struct{}{}
    }

    for _, tx := range a {
        if _, ok := remove[tx.hash]; !ok {
            keep = append(keep, tx)
        }
    }

    return keep
}

// TxByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxByNonceFast TransactionsFast

func (s TxByNonceFast) Len() int           { return len(s) }
func (s TxByNonceFast) Less(i, j int) bool { return s[i].Nonce() < s[j].Nonce() }
func (s TxByNonceFast) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// copyAddressPtr copies an address.
// func copyAddressPtr(a *common.Address) *common.Address {
//	if a == nil {
//		return nil
//	}
//	cpy := *a
//	return &cpy
// }
