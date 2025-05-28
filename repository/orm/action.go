package orm

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Action struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Maker        string
	Token        string
	Pair         string
	Action       string
	TxHash       string
	Creator      string
	Block        uint64
	BlockAt      time.Time
	Token0Amount decimal.Decimal `gorm:"-"`
	Token1Amount decimal.Decimal `gorm:"-"`
	CreatedAt    time.Time       `gorm:"autoCreateTime"`
}

func (a *Action) TableName() string {
	return "action"
}

func (a *Action) Equal(a2 *Action) bool {
	if a.Maker != a2.Maker {
		return false
	}
	if a.Token != a2.Token {
		return false
	}
	if a.Pair != a2.Pair {
		return false
	}
	if a.Action != a2.Action {
		return false
	}
	if a.TxHash != a2.TxHash {
		return false
	}
	if a.Creator != a2.Creator {
		return false
	}
	return true
}
