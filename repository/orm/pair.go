package orm

import (
	"github.com/shopspring/decimal"
	"time"
)

type Pair struct {
	Name      string
	Address   string
	Token0    string
	Token1    string
	ChainId   int
	Reserve0  decimal.Decimal
	Reserve1  decimal.Decimal
	Block     uint64
	BlockAt   time.Time
	Program   string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (p *Pair) TableName() string {
	return "pair"
}

func (p *Pair) Equal(p2 *Pair) bool {
	if p.Name != p2.Name {
		return false
	}
	if p.Address != p2.Address {
		return false
	}
	if p.Token0 != p2.Token0 {
		return false
	}
	if p.Token1 != p2.Token1 {
		return false
	}
	if p.ChainId != p2.ChainId {
		return false
	}
	if !p.Reserve0.Equal(p2.Reserve0) {
		return false
	}
	if !p.Reserve1.Equal(p2.Reserve1) {
		return false
	}
	return true
}
