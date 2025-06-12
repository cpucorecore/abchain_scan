package event_parser

import (
	"abchain_scan/parser/event_parser/event"
	"abchain_scan/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type SyncEventParser struct {
	PoolEventParser
}

func (o *SyncEventParser) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	input, err := o.ethLogUnpacker.Unpack(ethLog)
	if err != nil {
		return nil, err
	}

	e := &event.SyncEvent{
		EventCommon: types.EventCommonFromEthLog(ethLog),
		Amount0Wei:  input[0].(*big.Int),
		Amount1Wei:  input[1].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address: ethLog.Address,
	}

	e.PossibleProtocolIds = o.PossibleProtocolIds

	return e, nil
}
