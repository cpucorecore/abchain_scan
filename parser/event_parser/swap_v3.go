package event_parser

import (
	"abchain_scan/parser/event_parser/event"
	"abchain_scan/types"
	"errors"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

var (
	errAmount0Zero = errors.New("amount0 is zero")
	errAmount1Zero = errors.New("amount1 is zero")
)

type SwapEventParserV3 struct {
	PoolEventParser
}

func (o *SwapEventParserV3) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	input, err := o.ethLogUnpacker.Unpack(ethLog)
	if err != nil {
		return nil, err
	}

	e := &event.SwapEventV3{
		EventCommon: types.EventCommonFromEthLog(ethLog),
		Amount0Wei:  input[0].(*big.Int),
		Amount1Wei:  input[1].(*big.Int),
	}

	if e.Amount0Wei.Sign() == 0 {
		return nil, errAmount0Zero
	}

	if e.Amount1Wei.Sign() == 0 {
		return nil, errAmount1Zero
	}

	e.Pair = &types.Pair{
		Address: ethLog.Address,
	}

	e.PossibleProtocolIds = o.PossibleProtocolIds

	return e, nil
}
