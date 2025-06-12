package event_parser

import (
	"abchain_scan/repository/orm"
	"abchain_scan/service"
	"abchain_scan/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestMint_Aerodrome(t *testing.T) {
	// https://basescan.org/tx/0xff85824c89b77fb78641d11d20738817dbc7fdd0dbad9e791b4e8b2ad8f1a4e7#eventlog#72
	txHash := "0xff85824c89b77fb78641d11d20738817dbc7fdd0dbad9e791b4e8b2ad8f1a4e7"
	logIndex := 6
	pairAddress := "0xC09F68906B1DC60F1BA5771Ec6625cA947031Aaf"
	token0Address := "0x1E50309675d5C41D38Ba14133B4DB5b9f44FfBCd"
	block := uint64(30217994)
	BlockIndex := uint(186)
	TxIndex := uint(72)
	expectAmt0Wei, _ := decimal.NewFromString("10000000000000000000000")
	expectAmt1Wei, _ := decimal.NewFromString("4000000000000000000")
	program := types.ProtocolNameAerodrome

	tc := service.GetTestContext()
	ethLog := tc.GetEthLog(txHash, logIndex)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(service.MockNativeTokenPrice)

	token0Wei := decimal.NewFromBigInt(big.NewInt(1), int32(pairWrap.Pair.Token0Core.Decimals))
	expectAmt0 := expectAmt0Wei.Div(token0Wei)
	expectAmt1 := expectAmt1Wei.Div(service.Wei18)
	expectTx := &orm.Tx{
		TxHash:        txHash,
		Event:         types.Add,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Token0Address: token0Address,
		Token1Address: types.WETH,
		Block:         block,
		BlockIndex:    BlockIndex,
		TxIndex:       TxIndex,
		PairAddress:   pairAddress,
		Program:       program,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}

func TestMint_UniswapV2(t *testing.T) {
	// https://basescan.org/tx/0x2a925551ba86be62e96480791f1b152a6c9f542315c2c45643238e77be990b97#eventlog#386
	txHash := "0x2a925551ba86be62e96480791f1b152a6c9f542315c2c45643238e77be990b97"
	logIndex := 9
	pairAddress := "0xeD293D73563595E40c3354860721403f9E015CE4"
	token0Address := "0xF5D420881186855451eD4Fb4b6EF7B7E4f484F7D"
	block := uint64(30251018)
	BlockIndex := uint(117)
	TxIndex := uint(386)
	expectAmt0Wei, _ := decimal.NewFromString("100000000000000000000000000")
	expectAmt1Wei, _ := decimal.NewFromString("10000000000000")
	program := types.ProtocolNameNewSwap

	tc := service.GetTestContext()
	ethLog := tc.GetEthLog(txHash, logIndex)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(service.MockNativeTokenPrice)

	token0Wei := decimal.NewFromBigInt(big.NewInt(1), int32(pairWrap.Pair.Token0Core.Decimals))
	expectAmt0 := expectAmt0Wei.Div(token0Wei)
	expectAmt1 := expectAmt1Wei.Div(service.Wei18)
	expectTx := &orm.Tx{
		TxHash:        txHash,
		Event:         types.Add,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Token0Address: token0Address,
		Token1Address: types.WETH,
		Block:         block,
		BlockIndex:    BlockIndex,
		TxIndex:       TxIndex,
		PairAddress:   pairAddress,
		Program:       program,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}

func TestMint_PancakeV2(t *testing.T) {
	// https://basescan.org/tx/0xc5d3eae38ca43b8cb8beea8a07fb5b4ee91a805c5ff0569305a83a83a822e16f#eventlog#122
	txHash := "0xc5d3eae38ca43b8cb8beea8a07fb5b4ee91a805c5ff0569305a83a83a822e16f"
	logIndex := 8
	pairAddress := "0xabd5Ba8F0945d89E48077947037455bD40f475FC"
	token0Address := "0x357f9404356970F6B1C7208cf966abE3e505BA60"
	block := uint64(30230565)
	BlockIndex := uint(45)
	TxIndex := uint(122)
	expectAmt0Wei, _ := decimal.NewFromString("10000000000000000")
	expectAmt1Wei, _ := decimal.NewFromString("3500000000000000")
	program := types.ProtocolNamePancakeV2

	tc := service.GetTestContext()
	ethLog := tc.GetEthLog(txHash, logIndex)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(service.MockNativeTokenPrice)

	token0Wei := decimal.NewFromBigInt(big.NewInt(1), int32(pairWrap.Pair.Token0Core.Decimals))
	expectAmt0 := expectAmt0Wei.Div(token0Wei)
	expectAmt1 := expectAmt1Wei.Div(service.Wei18)
	expectTx := &orm.Tx{
		TxHash:        txHash,
		Event:         types.Add,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Token0Address: token0Address,
		Token1Address: types.WETH,
		Block:         block,
		BlockIndex:    BlockIndex,
		TxIndex:       TxIndex,
		PairAddress:   pairAddress,
		Program:       program,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}
