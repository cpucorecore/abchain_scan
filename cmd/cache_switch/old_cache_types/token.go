package old_cache_types

import (
	"abchain_scan/types"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

type Token struct {
	Address        common.Address
	Creator        common.Address
	Name           string
	Symbol         string
	Decimals       int16
	TotalSupply    decimal.Decimal
	BlockNumber    *big.Int
	BlockTime      time.Time
	Filtered       bool
	FilteredReason int
	Program        string
	MainPair       common.Address
}

func (t *Token) MarshalBinary() ([]byte, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (t *Token) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &t)
}

func (t *Token) ToNewToken() *types.Token {
	token := &types.Token{
		Address:     t.Address,
		Creator:     t.Creator,
		Name:        t.Name,
		Symbol:      t.Symbol,
		Decimals:    int8(t.Decimals),
		TotalSupply: t.TotalSupply,
		BlockTime:   t.BlockTime,
		Program:     t.Program,
		Filtered:    t.Filtered,
	}

	if t.BlockNumber != nil {
		token.BlockNumber = t.BlockNumber.Uint64()
	}

	return token
}
