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

func TestSwap_Aerodrome(t *testing.T) {
	// https://basescan.org/tx/0xc28b99ec1b68f67e08a8f519c7c972b87b5499b28ed0789cf48024eef275ca9b#eventlog#513
	txHash := "0xc28b99ec1b68f67e08a8f519c7c972b87b5499b28ed0789cf48024eef275ca9b"
	logIndex := 6
	pairAddress := "0xC09F68906B1DC60F1BA5771Ec6625cA947031Aaf"
	token0Address := "0x1E50309675d5C41D38Ba14133B4DB5b9f44FfBCd"
	block := uint64(30252205)
	BlockIndex := uint(120)
	TxIndex := uint(513)
	expectAmt0Wei, _ := decimal.NewFromString("20238")
	expectAmt1Wei, _ := decimal.NewFromString("199700000000000")
	eventName := types.Buy
	program := types.ProtocolNameAerodrome

	tc := service.GetTestContext()
	ethLog := tc.GetEthLog(txHash, logIndex)

	t.Log(ethLog.Topics[0].String())
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
		Event:         eventName,
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

func TestSwap_UniswapV2(t *testing.T) {
	// https://basescan.org/tx/0xaf9961fd1e289c02d0a70fe77e30ff8e47968e339f83b0d7c25c11fa72994fe8#eventlog#6
	txHash := "0xaf9961fd1e289c02d0a70fe77e30ff8e47968e339f83b0d7c25c11fa72994fe8"
	logIndex := 3
	pairAddress := "0xeD293D73563595E40c3354860721403f9E015CE4"
	token0Address := "0xF5D420881186855451eD4Fb4b6EF7B7E4f484F7D"
	block := uint64(30251020)
	BlockIndex := uint(4)
	TxIndex := uint(6)
	expectAmt0Wei, _ := decimal.NewFromString("28358176223964627199916808")
	expectAmt1Wei, _ := decimal.NewFromString("5000000000000")
	eventName := types.Buy
	program := types.ProtocolNameUniswapV2

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
		Event:         eventName,
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

func TestSwap_PancakeV2(t *testing.T) {
	// TODO
}
