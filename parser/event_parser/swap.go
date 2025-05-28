package event_parser

import (
	"abchain_scan/parser/event_parser/event"
	"abchain_scan/types"
	"errors"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

var (
	errAmountInZero  = errors.New("amount in is zero")
	errAmountOutZero = errors.New("amount out is zero")
)

type SwapEventParser struct {
	PoolEventParser
}

func (o *SwapEventParser) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	eventInput, err := o.ethLogUnpacker.Unpack(ethLog)
	if err != nil {
		return nil, err
	}

	e := &event.SwapEvent{
		EventCommon:   types.EventCommonFromEthLog(ethLog),
		Amount0InWei:  eventInput[0].(*big.Int),
		Amount1InWei:  eventInput[1].(*big.Int),
		Amount0OutWei: eventInput[2].(*big.Int),
		Amount1OutWei: eventInput[3].(*big.Int),
	}

	if e.Amount0InWei.Sign() == 0 && e.Amount1InWei.Sign() == 0 {
		return nil, errAmountInZero
	}

	if e.Amount0OutWei.Sign() == 0 && e.Amount1OutWei.Sign() == 0 {
		return nil, errAmountOutZero
	}

	e.Pair = &types.Pair{
		Address: ethLog.Address,
	}

	e.PossibleProtocolIds = o.PossibleProtocolIds

	return e, nil
}
