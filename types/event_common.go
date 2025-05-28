package types

import (
	"abchain_scan/repository/orm"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"time"
)

type EventCommon struct {
	Pair                *Pair
	ContractAddress     common.Address
	BlockNumber         uint64
	BlockTime           time.Time
	TxHash              common.Hash
	Maker               common.Address
	TxIndex             uint
	LogIndex            uint
	PossibleProtocolIds []int
}

var _ Event = &EventCommon{}

func (e *EventCommon) GetProtocolId() int {
	return e.Pair.ProtocolId
}

func (e *EventCommon) GetPossibleProtocolIds() []int {
	return e.PossibleProtocolIds
}

func (e *EventCommon) SetPossibleProtocolIds(possibleProtocolIds []int) {
	e.PossibleProtocolIds = possibleProtocolIds
}

func (e *EventCommon) CanGetPair() bool {
	return false
}

func (e *EventCommon) GetPair() *Pair {
	return nil
}

func (e *EventCommon) GetPairAddress() common.Address {
	return e.Pair.Address
}

func (e *EventCommon) SetPair(pair *Pair) {
	e.Pair = pair
}

func (e *EventCommon) CanGetTx() bool {
	return false
}

func (e *EventCommon) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	return nil
}

func (e *EventCommon) CanGetPoolUpdate() bool {
	return false
}

func (e *EventCommon) GetPoolUpdate() *PoolUpdate {
	return nil
}

func (e *EventCommon) CanGetPoolUpdateParameter() bool {
	return false
}

func (e *EventCommon) GetPoolUpdateParameter() *PoolUpdateParameter {
	return nil
}

func (e *EventCommon) LinkEvent(event Event) {
}

func (e *EventCommon) IsCreatePair() bool {
	return false
}

func (e *EventCommon) IsMint() bool {
	return false
}

func (e *EventCommon) GetMintAmount() (decimal.Decimal, decimal.Decimal) {
	return decimal.Zero, decimal.Zero
}

func (e *EventCommon) SetMaker(maker common.Address) {
	e.Maker = maker
}

func (e *EventCommon) SetBlockTime(blockTime time.Time) {
	e.BlockTime = blockTime
}

func EventCommonFromEthLog(ethLog *ethtypes.Log) *EventCommon {
	return &EventCommon{
		ContractAddress: ethLog.Address,
		BlockNumber:     ethLog.BlockNumber,
		TxHash:          ethLog.TxHash,
		TxIndex:         ethLog.TxIndex,
		LogIndex:        ethLog.Index,
	}
}
