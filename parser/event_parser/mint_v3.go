package event_parser

import (
	"abchain_scan/parser/event_parser/event"
	"abchain_scan/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type MintEventParserV3 struct {
	PoolEventParser
}

func (o *MintEventParserV3) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	input, err := o.ethLogUnpacker.Unpack(ethLog)
	if err != nil {
		return nil, err
	}

	e := &event.MintEvent{
		EventCommon: types.EventCommonFromEthLog(ethLog),
		Amount0Wei:  input[2].(*big.Int),
		Amount1Wei:  input[3].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address: ethLog.Address,
	}

	e.PossibleProtocolIds = o.PossibleProtocolIds

	return e, nil
}
