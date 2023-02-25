// Copyright 2021 The go-ethereum Authors
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
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// txJSON is the JSON representation of transactions.
type txFastJSON struct {
	Type hexutil.Uint64 `json:"type"`

	// Common transaction fields:
	Nonce                *hexutil.Uint64 `json:"nonce"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	Value                *hexutil.Big    `json:"value"`
	Data                 *hexutil.Bytes  `json:"input"`
	R                    *hexutil.Big    `json:"r"`
	To                   *string `json:"to"`
	From                 *string `json:"from,omitempty"`

	// Access list transaction fields:
	ChainID    *hexutil.Uint64 `json:"chainId,omitempty"`
	BlockNumber *hexutil.Uint64 `json:"blockNumber,omitempty"`

	// Only used for encoding:
	Hash common.Hash `json:"hash"`
}

// MarshalJSON marshals as JSON with a hash.
func (t *TransactionFast) MarshalJSON() ([]byte, error) {
	return nil, errors.New("It is not posible to marshal this custom structure")
}

// UnmarshalJSON unmarshals from JSON.
func (t *TransactionFast) UnmarshalJSON(input []byte) error {
	var dec txFastJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	// Decode / verify fields according to transaction type.
	switch dec.Type {
	case LegacyTxType:

		t.to = ""
		if dec.To != nil {
			t.to = *dec.To
		}
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		t.nonce = uint64(*dec.Nonce)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		t.gasPrice = (*big.Int)(dec.GasPrice)
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		t.gas = uint64(*dec.Gas)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		t.value = (*big.Int)(dec.Value)
		if dec.Data == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		t.data = *dec.Data
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		t.r = (*big.Int)(dec.R)

		t.from = ""
		if dec.From != nil {
			t.from = *dec.From
		}
		t.BlockNumber = 0
		if dec.BlockNumber != nil {
			t.BlockNumber = uint64(*dec.BlockNumber)
		}

	//case AccessListTxType:
	//	var itx AccessListTx
	//	// Access list is optional for now.
	//	if dec.AccessList != nil {
	//		itx.AccessList = *dec.AccessList
	//	}
	//	if dec.ChainID == nil {
	//		return errors.New("missing required field 'chainId' in transaction")
	//	}
	//	itx.ChainID = (*big.Int)(dec.ChainID)
	//	if dec.To != nil {
	//		itx.To = dec.To
	//	}
	//	if dec.Nonce == nil {
	//		return errors.New("missing required field 'nonce' in transaction")
	//	}
	//	itx.Nonce = uint64(*dec.Nonce)
	//	if dec.GasPrice == nil {
	//		return errors.New("missing required field 'gasPrice' in transaction")
	//	}
	//	itx.GasPrice = (*big.Int)(dec.GasPrice)
	//	if dec.Gas == nil {
	//		return errors.New("missing required field 'gas' in transaction")
	//	}
	//	itx.Gas = uint64(*dec.Gas)
	//	if dec.Value == nil {
	//		return errors.New("missing required field 'value' in transaction")
	//	}
	//	itx.Value = (*big.Int)(dec.Value)
	//	if dec.Data == nil {
	//		return errors.New("missing required field 'input' in transaction")
	//	}
	//	itx.Data = *dec.Data
	//	if dec.V == nil {
	//		return errors.New("missing required field 'v' in transaction")
	//	}
	//	itx.V = (*big.Int)(dec.V)
	//	if dec.R == nil {
	//		return errors.New("missing required field 'r' in transaction")
	//	}
	//	itx.R = (*big.Int)(dec.R)
	//	if dec.S == nil {
	//		return errors.New("missing required field 's' in transaction")
	//	}
	//	itx.S = (*big.Int)(dec.S)
	//	withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
	//	if withSignature {
	//		if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
	//			return err
	//		}
	//	}

	//case DynamicFeeTxType:
	//	var itx DynamicFeeTx
	//	inner = &itx
	//	// Access list is optional for now.
	//	if dec.AccessList != nil {
	//		itx.AccessList = *dec.AccessList
	//	}
	//	if dec.ChainID == nil {
	//		return errors.New("missing required field 'chainId' in transaction")
	//	}
	//	itx.ChainID = (*big.Int)(dec.ChainID)
	//	if dec.To != nil {
	//		itx.To = dec.To
	//	}
	//	if dec.Nonce == nil {
	//		return errors.New("missing required field 'nonce' in transaction")
	//	}
	//	itx.Nonce = uint64(*dec.Nonce)
	//	if dec.MaxPriorityFeePerGas == nil {
	//		return errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
	//	}
	//	itx.GasTipCap = (*big.Int)(dec.MaxPriorityFeePerGas)
	//	if dec.MaxFeePerGas == nil {
	//		return errors.New("missing required field 'maxFeePerGas' for txdata")
	//	}
	//	itx.GasFeeCap = (*big.Int)(dec.MaxFeePerGas)
	//	if dec.Gas == nil {
	//		return errors.New("missing required field 'gas' for txdata")
	//	}
	//	itx.Gas = uint64(*dec.Gas)
	//	if dec.Value == nil {
	//		return errors.New("missing required field 'value' in transaction")
	//	}
	//	itx.Value = (*big.Int)(dec.Value)
	//	if dec.Data == nil {
	//		return errors.New("missing required field 'input' in transaction")
	//	}
	//	itx.Data = *dec.Data
	//	if dec.V == nil {
	//		return errors.New("missing required field 'v' in transaction")
	//	}
	//	itx.V = (*big.Int)(dec.V)
	//	if dec.R == nil {
	//		return errors.New("missing required field 'r' in transaction")
	//	}
	//	itx.R = (*big.Int)(dec.R)
	//	if dec.S == nil {
	//		return errors.New("missing required field 's' in transaction")
	//	}
	//	itx.S = (*big.Int)(dec.S)
	//	withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
	//	if withSignature {
	//		if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
	//			return err
	//		}
	//	}

	default:
		return ErrTxTypeNotSupported
	}

	// Now set the inner transaction.
	t.time = time.Now()

	// TODO: check hash here?
	return nil
}
