package event_parser

import (
	"abchain_scan/service"
	"abchain_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestSync_Aerodrome(t *testing.T) {
	// https://basescan.org/tx/0xeb50cf26b45a8ca72d9343dab433f5eeecaeb3eabc2716a7dd80ddef966948b1#eventlog#212
	txHash := "0xeb50cf26b45a8ca72d9343dab433f5eeecaeb3eabc2716a7dd80ddef966948b1"
	logIndex := 5
	pairAddress := "0xC09F68906B1DC60F1BA5771Ec6625cA947031Aaf"
	token0Address := "0x1E50309675d5C41D38Ba14133B4DB5b9f44FfBCd"
	LogIndex := uint(212)
	expectAmt0Wei, _ := decimal.NewFromString("20239")
	expectAmt1Wei, _ := decimal.NewFromString("50")
	program := types.ProtocolNameAerodrome

	tc := service.GetTestContext()
	ethLog := tc.GetEthLog(txHash, logIndex)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	require.False(t, event.CanGetTx())
	require.True(t, event.CanGetPoolUpdate())
	poolUpdate := event.GetPoolUpdate()
	token0Wei := decimal.NewFromBigInt(big.NewInt(1), int32(pairWrap.Pair.Token0Core.Decimals))
	expectAmt0 := expectAmt0Wei.Div(token0Wei)
	expectAmt1 := expectAmt1Wei.Div(service.Wei18)
	expectPoolUpdate := &types.PoolUpdate{
		Program:       program,
		LogIndex:      LogIndex,
		Address:       common.HexToAddress(pairAddress),
		Token0Address: common.HexToAddress(token0Address),
		Token1Address: types.WETHAddress,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
	}
	require.True(t, expectPoolUpdate.Equal(poolUpdate), "expect: %v, actual: %v", expectPoolUpdate, poolUpdate)
}

func TestSync_UniswapV2(t *testing.T) {
	// https://basescan.org/tx/0x2a925551ba86be62e96480791f1b152a6c9f542315c2c45643238e77be990b97#eventlog#385
	txHash := "0x2a925551ba86be62e96480791f1b152a6c9f542315c2c45643238e77be990b97"
	logIndex := 8
	pairAddress := "0xeD293D73563595E40c3354860721403f9E015CE4"
	token0Address := "0xF5D420881186855451eD4Fb4b6EF7B7E4f484F7D"
	LogIndex := uint(385)
	expectAmt0Wei, _ := decimal.NewFromString("100000000000000000000000000")
	expectAmt1Wei, _ := decimal.NewFromString("10000000000000")
	program := types.ProtocolNameUniswapV2

	tc := service.GetTestContext()
	ethLog := tc.GetEthLog(txHash, logIndex)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	require.False(t, event.CanGetTx())
	require.True(t, event.CanGetPoolUpdate())
	poolUpdate := event.GetPoolUpdate()
	token0Wei := decimal.NewFromBigInt(big.NewInt(1), int32(pairWrap.Pair.Token0Core.Decimals))
	expectAmt0 := expectAmt0Wei.Div(token0Wei)
	expectAmt1 := expectAmt1Wei.Div(service.Wei18)
	expectPoolUpdate := &types.PoolUpdate{
		Program:       program,
		LogIndex:      LogIndex,
		Address:       common.HexToAddress(pairAddress),
		Token0Address: common.HexToAddress(token0Address),
		Token1Address: types.WETHAddress,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
	}
	require.True(t, expectPoolUpdate.Equal(poolUpdate), "expect: %v, actual: %v", expectPoolUpdate, poolUpdate)
}

func TestSync_PancakeV2(t *testing.T) {
	// https://basescan.org/tx/0xc5d3eae38ca43b8cb8beea8a07fb5b4ee91a805c5ff0569305a83a83a822e16f#eventlog#121
	txHash := "0xc5d3eae38ca43b8cb8beea8a07fb5b4ee91a805c5ff0569305a83a83a822e16f"
	logIndex := 7
	pairAddress := "0xabd5Ba8F0945d89E48077947037455bD40f475FC"
	token0Address := "0x357f9404356970F6B1C7208cf966abE3e505BA60"
	LogIndex := uint(121)
	expectAmt0Wei, _ := decimal.NewFromString("10000000000000000")
	expectAmt1Wei, _ := decimal.NewFromString("3500000000000000")
	program := types.ProtocolNamePancakeV2

	tc := service.GetTestContext()
	ethLog := tc.GetEthLog(txHash, logIndex)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	require.False(t, event.CanGetTx())
	require.True(t, event.CanGetPoolUpdate())
	poolUpdate := event.GetPoolUpdate()
	token0Wei := decimal.NewFromBigInt(big.NewInt(1), int32(pairWrap.Pair.Token0Core.Decimals))
	expectAmt0 := expectAmt0Wei.Div(token0Wei)
	expectAmt1 := expectAmt1Wei.Div(service.Wei18)
	expectPoolUpdate := &types.PoolUpdate{
		Program:       program,
		LogIndex:      LogIndex,
		Address:       common.HexToAddress(pairAddress),
		Token0Address: common.HexToAddress(token0Address),
		Token1Address: types.WETHAddress,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
	}
	require.True(t, expectPoolUpdate.Equal(poolUpdate), "expect: %v, actual: %v", expectPoolUpdate, poolUpdate)
}
