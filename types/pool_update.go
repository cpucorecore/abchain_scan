package types

import (
	"abchain_scan/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type PoolUpdate struct {
	Program       string // TODO remove
	LogIndex      uint
	Address       common.Address
	Token0Address common.Address
	Token1Address common.Address
	Token0Amount  decimal.Decimal
	Token1Amount  decimal.Decimal
}

func (u *PoolUpdate) Equal(tx *PoolUpdate) bool {
	if u.Program != tx.Program {
		return false
	}
	if u.LogIndex != tx.LogIndex {
		return false
	}
	if u.Address != tx.Address {
		return false
	}
	if u.Token0Address != tx.Token0Address {
		return false
	}
	if u.Token1Address != tx.Token1Address {
		return false
	}
	if !util.DecimalEqual(u.Token0Amount, tx.Token0Amount) {
		return false
	}
	if !util.DecimalEqual(u.Token1Amount, tx.Token1Amount) {
		return false
	}
	return true
}

type PoolUpdateParameter struct {
	BlockNumber   uint64
	PairAddress   common.Address
	Token0Address common.Address
	Token1Address common.Address
}
