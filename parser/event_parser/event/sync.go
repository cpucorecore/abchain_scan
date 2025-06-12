package event

import (
	"abchain_scan/types"
	"math/big"
)

type SyncEvent struct {
	*types.EventCommon
	Amount0Wei *big.Int
	Amount1Wei *big.Int
}

func (e *SyncEvent) CanGetPoolUpdate() bool {
	return true
}

func (e *SyncEvent) GetPoolUpdate() *types.PoolUpdate {
	pu := &types.PoolUpdate{
		Program:       types.GetProtocolName(e.GetProtocolId()),
		LogIndex:      e.LogIndex,
		Address:       e.ContractAddress,
		Token0Address: e.Pair.Token0Core.Address,
		Token1Address: e.Pair.Token1Core.Address,
	}

	pu.Token0Amount, pu.Token1Amount = ParseAmountsByPair(e.Amount0Wei, e.Amount1Wei, e.Pair)

	return pu
}

var _ types.Event = (*SyncEvent)(nil)
