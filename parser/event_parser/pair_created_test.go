package event_parser

import (
	"abchain_scan/service"
	"abchain_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPairCreated_Aerodrome(t *testing.T) {
	// https://basescan.org/tx/0xff85824c89b77fb78641d11d20738817dbc7fdd0dbad9e791b4e8b2ad8f1a4e7#eventlog#66
	tc := service.GetTestContext()
	ethLog := tc.GetEthLog("0xff85824c89b77fb78641d11d20738817dbc7fdd0dbad9e791b4e8b2ad8f1a4e7", 0)
	blockTimestamp := tc.GetBlockTimestamp(ethLog.BlockNumber)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	event.SetBlockTime(time.Unix(int64(blockTimestamp), 0))
	require.True(t, event.CanGetPair())
	pair := event.GetPair()
	pairWrap := tc.PairService.GetPairTokens(pair)
	event.SetPair(pairWrap.Pair)

	expectPair := &types.Pair{
		Address:        common.HexToAddress("0xC09F68906B1DC60F1BA5771Ec6625cA947031Aaf"),
		TokensReversed: false,
		Token0Core: &types.TokenCore{
			Address:  common.HexToAddress("0x1E50309675d5C41D38Ba14133B4DB5b9f44FfBCd"),
			Symbol:   "NOHUMAIN",
			Decimals: 18,
		},
		Token1Core: &types.TokenCore{
			Address:  types.WETHAddress,
			Symbol:   "WETH",
			Decimals: 18,
		},
		Block:      30217994,
		BlockAt:    time.Unix(int64(blockTimestamp), 0),
		ProtocolId: types.ProtocolIdAerodrome,
		Filtered:   false,
		FilterCode: 0,
	}

	require.True(t, pairWrap.Pair.Equal(expectPair), "expect: %v, actual: %v", expectPair, pairWrap.Pair)
}

func TestPairCreated_UniswapV2(t *testing.T) {
	// https://basescan.org/tx/0x2a925551ba86be62e96480791f1b152a6c9f542315c2c45643238e77be990b97#eventlog#377
	tc := service.GetTestContext()
	ethLog := tc.GetEthLog("0x2a925551ba86be62e96480791f1b152a6c9f542315c2c45643238e77be990b97", 0)
	blockTimestamp := tc.GetBlockTimestamp(ethLog.BlockNumber)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	event.SetBlockTime(time.Unix(int64(blockTimestamp), 0))
	require.True(t, event.CanGetPair())
	pair := event.GetPair()
	pairWrap := tc.PairService.GetPairTokens(pair)
	event.SetPair(pairWrap.Pair)

	expectPair := &types.Pair{
		Address:        common.HexToAddress("0xeD293D73563595E40c3354860721403f9E015CE4"),
		TokensReversed: true,
		Token0Core: &types.TokenCore{
			Address:  common.HexToAddress("0xF5D420881186855451eD4Fb4b6EF7B7E4f484F7D"),
			Symbol:   "AE JasmyCoin",
			Decimals: 18,
		},
		Token1Core: &types.TokenCore{
			Address:  types.WETHAddress,
			Symbol:   "WETH",
			Decimals: 18,
		},
		Block:      30251018,
		BlockAt:    time.Unix(int64(blockTimestamp), 0),
		ProtocolId: types.ProtocolIdUniswapV2,
		Filtered:   false,
		FilterCode: 0,
	}

	require.True(t, pairWrap.Pair.Equal(expectPair), "expect: %v, actual: %v", expectPair, pairWrap.Pair)
}

func TestPairCreated_PancakeV2(t *testing.T) {
	// https://basescan.org/tx/0xc5d3eae38ca43b8cb8beea8a07fb5b4ee91a805c5ff0569305a83a83a822e16f#eventlog#114
	tc := service.GetTestContext()
	ethLog := tc.GetEthLog("0xc5d3eae38ca43b8cb8beea8a07fb5b4ee91a805c5ff0569305a83a83a822e16f", 0)
	blockTimestamp := tc.GetBlockTimestamp(ethLog.BlockNumber)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	event.SetBlockTime(time.Unix(int64(blockTimestamp), 0))
	require.True(t, event.CanGetPair())
	pair := event.GetPair()
	pairWrap := tc.PairService.GetPairTokens(pair)
	event.SetPair(pairWrap.Pair)

	expectPair := &types.Pair{
		Address:        common.HexToAddress("0xabd5Ba8F0945d89E48077947037455bD40f475FC"),
		TokensReversed: false,
		Token0Core: &types.TokenCore{
			Address:  common.HexToAddress("0x357f9404356970F6B1C7208cf966abE3e505BA60"),
			Symbol:   "Mutuum Presale",
			Decimals: 8,
		},
		Token1Core: &types.TokenCore{
			Address:  types.WETHAddress,
			Symbol:   "WETH",
			Decimals: 18,
		},
		Block:      30230565,
		BlockAt:    time.Unix(int64(blockTimestamp), 0),
		ProtocolId: types.ProtocolIdPancakeV2,
		Filtered:   false,
		FilterCode: 0,
	}

	require.True(t, pairWrap.Pair.Equal(expectPair), "expect: %v, actual: %v", expectPair, pairWrap.Pair)
}
