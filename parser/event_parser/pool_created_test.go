package event_parser

import (
	"abchain_scan/service"
	"abchain_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPoolCreated_UniswapV3(t *testing.T) {
	// https://basescan.org/tx/0x91db85460d929bfc664779b1bf7fc23ea47436f84bdd1782e6466cb0bb2962ef#eventlog#797
	tc := service.GetTestContext()
	ethLog := tc.GetEthLog("0x91db85460d929bfc664779b1bf7fc23ea47436f84bdd1782e6466cb0bb2962ef", 2)
	blockTimestamp := tc.GetBlockTimestamp(ethLog.BlockNumber)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	event.SetBlockTime(time.Unix(int64(blockTimestamp), 0))
	require.True(t, event.CanGetPair())
	pair := event.GetPair()
	pairWrap := tc.PairService.GetPairTokens(pair)
	event.SetPair(pairWrap.Pair)

	expectPair := &types.Pair{
		Address:        common.HexToAddress("0x2c93555C0150DA726957782e36A12D76D6851064"),
		TokensReversed: true,
		Token0Core: &types.TokenCore{
			Address:  common.HexToAddress("0x95f51e058AB0104211659548A728A849A47A0b07"),
			Symbol:   "fuck",
			Decimals: 18,
		},
		Token1Core: &types.TokenCore{
			Address:  types.WETHAddress,
			Symbol:   "WETH",
			Decimals: 18,
		},
		Block:      30253895,
		BlockAt:    time.Unix(int64(blockTimestamp), 0),
		ProtocolId: types.ProtocolIdUniswapV3,
		Filtered:   false,
		FilterCode: 0,
	}

	require.True(t, pairWrap.Pair.Equal(expectPair), "expect: %v, actual: %v", expectPair, pairWrap.Pair)
}

func TestPoolCreated_PancakeV3(t *testing.T) {
	// https://basescan.org/tx/0x0809776ecaeb651bd0f354e85814e7ed792709e5397ab777044a82a4a760d445#eventlog#291
	tc := service.GetTestContext()
	ethLog := tc.GetEthLog("0x0809776ecaeb651bd0f354e85814e7ed792709e5397ab777044a82a4a760d445", 0)
	blockTimestamp := tc.GetBlockTimestamp(ethLog.BlockNumber)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	event.SetBlockTime(time.Unix(int64(blockTimestamp), 0))
	require.True(t, event.CanGetPair())
	pair := event.GetPair()
	pairWrap := tc.PairService.GetPairTokens(pair)
	event.SetPair(pairWrap.Pair)

	expectPair := &types.Pair{
		Address:        common.HexToAddress("0x5845A51630AFab7C68556CF57a7b6827Bd94d434"),
		TokensReversed: true,
		Token0Core: &types.TokenCore{
			Address:  common.HexToAddress("0x858A6594f86fafb10dC0dEdC588C7CBb8E795129"),
			Symbol:   "JESUS",
			Decimals: 18,
		},
		Token1Core: &types.TokenCore{
			Address:  types.WETHAddress,
			Symbol:   "WETH",
			Decimals: 18,
		},
		Block:      30250940,
		BlockAt:    time.Unix(int64(blockTimestamp), 0),
		ProtocolId: types.ProtocolIdPancakeV3,
		Filtered:   false,
		FilterCode: 0,
	}

	require.True(t, pairWrap.Pair.Equal(expectPair), "expect: %v, actual: %v", expectPair, pairWrap.Pair)
}
