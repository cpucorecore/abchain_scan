package types

import (
	"abchain_scan/repository/orm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"time"
)

type Event interface {
	GetProtocolId() int
	GetPossibleProtocolIds() []int
	CanGetPair() bool
	GetPair() *Pair
	GetPairAddress() common.Address
	SetPair(pair *Pair)
	SetMaker(maker common.Address)
	SetBlockTime(blockTime time.Time)

	CanGetTx() bool
	GetTx(bnbPrice decimal.Decimal) *orm.Tx

	CanGetPoolUpdate() bool
	GetPoolUpdate() *PoolUpdate

	CanGetPoolUpdateParameter() bool
	GetPoolUpdateParameter() *PoolUpdateParameter

	LinkEvent(event Event)

	IsCreatePair() bool
	IsMint() bool
	GetMintAmount() (decimal.Decimal, decimal.Decimal)
}
