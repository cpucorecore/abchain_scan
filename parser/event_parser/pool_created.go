package event_parser

import (
	"abchain_scan/abi"
	"abchain_scan/parser/event_parser/event"
	"abchain_scan/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type PoolCreatedEventParser struct {
	FactoryEventParser
}

func (o *PoolCreatedEventParser) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	pair := &types.Pair{}

	_, ok := o.PossibleFactoryAddresses[ethLog.Address]
	if !ok {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeWrongFactory
		return nil, ErrWrongFactoryAddress
	}

	input, err := o.LogUnpacker.Unpack(ethLog)
	if err != nil {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeUnpackDataErr
		return nil, err
	}

	e := &event.PairCreatedEvent{
		EventCommon: types.EventCommonFromEthLog(ethLog),
	}

	pair.Address = input[1].(common.Address)
	pair.Token0Core = &types.TokenCore{
		Address: common.BytesToAddress(ethLog.Topics[1].Bytes()[12:]),
	}
	pair.Token1Core = &types.TokenCore{
		Address: common.BytesToAddress(ethLog.Topics[2].Bytes()[12:]),
	}
	pair.Block = ethLog.BlockNumber
	pair.ProtocolId = abi.FactoryAddress2ProtocolId[ethLog.Address]

	pair.FilterByToken0AndToken1()

	e.Pair = pair

	return e, nil
}
