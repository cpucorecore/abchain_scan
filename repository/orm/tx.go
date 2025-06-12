package orm

import (
	"abchain_scan/util"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Tx struct {
	Id            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;readonly"`
	TxHash        string
	Event         string
	Token0Amount  decimal.Decimal
	Token1Amount  decimal.Decimal
	Maker         string
	Token0Address string
	Token1Address string
	AmountUsd     decimal.Decimal
	PriceUsd      decimal.Decimal
	Block         uint64
	BlockAt       time.Time
	BlockIndex    uint
	TxIndex       uint
	PairAddress   string
	Program       string
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (t *Tx) Equal(tx *Tx) bool {
	if t.TxHash != tx.TxHash {
		return false
	}
	if t.Event != tx.Event {
		return false
	}
	//if t.Maker != tx.Maker {
	//	return false
	//}
	if t.Token0Address != tx.Token0Address {
		return false
	}
	if t.Token1Address != tx.Token1Address {
		return false
	}
	if t.Block != tx.Block {
		return false
	}
	if t.BlockIndex != tx.BlockIndex {
		return false
	}
	if t.TxIndex != tx.TxIndex {
		return false
	}
	if t.PairAddress != tx.PairAddress {
		return false
	}
	if t.Program != tx.Program {
		return false
	}
	if !util.DecimalEqual(tx.Token0Amount, t.Token0Amount) {
		return false
	}
	if !util.DecimalEqual(tx.Token1Amount, t.Token1Amount) {
		return false
	}

	// ignore priceUsd and amountUSd
	return true
}

func (t *Tx) TableName() string {
	return "tx"
}
