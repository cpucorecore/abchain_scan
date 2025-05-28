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

func TestBurn_UniswapV3(t *testing.T) {
	// https://basescan.org/tx/0x20cc88585f696d4f297691dd0bb6a8948d4b5267e30d1a65130ed4e2ef028aab#eventlog#68
	txHash := "0x20cc88585f696d4f297691dd0bb6a8948d4b5267e30d1a65130ed4e2ef028aab"
	logIndex := 0
	pairAddress := "0xd0b53D9277642d899DF5C87A3966A349A798F224" // WETH/USDC
	token0Address := "0x4200000000000000000000000000000000000006"
	block := uint64(30295298)
	BlockIndex := uint(33)
	TxIndex := uint(68)
	expectAmt0Wei, _ := decimal.NewFromString("17481434677215331194")
	expectAmt1Wei, _ := decimal.NewFromString("63861519633")
	program := types.ProtocolNameUniswapV3

	tc := service.GetTestContext()
	ethLog := tc.GetEthLog(txHash, logIndex)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(service.MockNativeTokenPrice)

	token0Wei := decimal.NewFromBigInt(big.NewInt(1), int32(pairWrap.Pair.Token0Core.Decimals))
	expectAmt0 := expectAmt0Wei.Div(token0Wei)
	expectAmt1 := expectAmt1Wei.Div(service.Wei6)
	expectTx := &orm.Tx{
		TxHash:        txHash,
		Event:         types.Remove,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Token0Address: token0Address,
		Token1Address: types.USDC,
		Block:         block,
		BlockIndex:    BlockIndex,
		TxIndex:       TxIndex,
		PairAddress:   pairAddress,
		Program:       program,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}

func TestBurn_PancakeV3(t *testing.T) {
	// https://basescan.org/tx/0xc7dd1b0488281af87b51b1c6362b898170c759413ac6770f183c1bdf38899dbe#eventlog#501
	txHash := "0xc7dd1b0488281af87b51b1c6362b898170c759413ac6770f183c1bdf38899dbe"
	logIndex := 0
	pairAddress := "0xb94b22332ABf5f89877A14Cc88f2aBC48c34B3Df" // USDC/WBTC
	token0Address := "0xcbB7C0000aB88B473b1f5aFd9ef808440eed33Bf"
	block := uint64(30298851)
	BlockIndex := uint(137)
	TxIndex := uint(501)
	expectAmt0Wei, _ := decimal.NewFromString("35413068")
	expectAmt1Wei, _ := decimal.NewFromString("22557511932")
	program := types.ProtocolNamePancakeV3

	tc := service.GetTestContext()
	ethLog := tc.GetEthLog(txHash, logIndex)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(service.MockNativeTokenPrice)

	token0Wei := decimal.NewFromBigInt(big.NewInt(1), int32(pairWrap.Pair.Token0Core.Decimals))
	token1Wei := decimal.NewFromBigInt(big.NewInt(1), int32(pairWrap.Pair.Token1Core.Decimals))
	expectAmt0 := expectAmt0Wei.Div(token0Wei)
	expectAmt1 := expectAmt1Wei.Div(token1Wei)
	expectTx := &orm.Tx{
		TxHash:        txHash,
		Event:         types.Remove,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Token0Address: token0Address,
		Token1Address: types.USDC,
		Block:         block,
		BlockIndex:    BlockIndex,
		TxIndex:       TxIndex,
		PairAddress:   pairAddress,
		Program:       program,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}
