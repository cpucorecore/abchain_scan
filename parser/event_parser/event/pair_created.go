package event

import "abchain_scan/types"

type PairCreatedEvent struct {
	*types.EventCommon
	MintEvent types.Event
}

func (e *PairCreatedEvent) CanGetPair() bool {
	return true
}

func (e *PairCreatedEvent) GetPair() *types.Pair {
	if e.MintEvent != nil {
		e.Pair.Token0InitAmount, e.Pair.Token1InitAmount = e.MintEvent.GetMintAmount()
	}
	e.Pair.BlockAt = e.BlockTime
	return e.Pair
}

func (e *PairCreatedEvent) IsCreatePair() bool {
	return true
}

func (e *PairCreatedEvent) LinkEvent(event types.Event) { // for pair initial token0/token1 amount
	e.MintEvent = event
}

var _ types.Event = (*PairCreatedEvent)(nil)
