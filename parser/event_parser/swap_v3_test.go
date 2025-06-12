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

func TestSwap_UniswapV3(t *testing.T) {
	// https://basescan.org/tx/0xc95dca773778f44671cc11d9b229332862df9c80775320004bf57d6644742b0c#eventlog#319
	txHash := "0xc95dca773778f44671cc11d9b229332862df9c80775320004bf57d6644742b0c"
	logIndex := 15
	pairAddress := "0x56f0eB23116F893feA120d1348E06548Fbc21af4"
	token0Address := "0x6d52F6Df8aF765749Fa4Ee3112F04FDb5C48AB07"
	block := uint64(30255154)
	BlockIndex := uint(112)
	TxIndex := uint(319)
	expectAmt0Wei, _ := decimal.NewFromString("2005862078597804492687218")
	expectAmt1Wei, _ := decimal.NewFromString("200000000000000")
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
	expectAmt1 := expectAmt1Wei.Div(service.Wei18)
	expectTx := &orm.Tx{
		TxHash:        txHash,
		Event:         types.Buy,
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

func TestSwap_PancakeV3(t *testing.T) {
	// https://basescan.org/tx/0x8b81b4cb2e09825e6d5b1c405e593c4a2f5531fec3850d5cbc2dcbe8f4eb12c6#eventlog#814
	txHash := "0x8b81b4cb2e09825e6d5b1c405e593c4a2f5531fec3850d5cbc2dcbe8f4eb12c6"
	logIndex := 8
	pairAddress := "0xb4024c8eBd3364505ABae90e6fC608763DA4d57c"
	token0Address := "0x99d079545B51043B4904c8986c4cA8fD0e64945D"
	block := uint64(30193945)
	BlockIndex := uint(221)
	TxIndex := uint(814)
	expectAmt0Wei, _ := decimal.NewFromString("409605315313239")
	expectAmt1Wei, _ := decimal.NewFromString("887371421897")
	program := types.ProtocolNamePancakeV3

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
		Event:         types.Buy,
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
