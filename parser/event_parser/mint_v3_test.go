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

func TestMint_UniswapV3(t *testing.T) {
	// https://basescan.org/tx/0x91db85460d929bfc664779b1bf7fc23ea47436f84bdd1782e6466cb0bb2962ef#eventlog#800
	txHash := "0x91db85460d929bfc664779b1bf7fc23ea47436f84bdd1782e6466cb0bb2962ef"
	logIndex := 5
	pairAddress := "0x2c93555C0150DA726957782e36A12D76D6851064"
	token0Address := "0x95f51e058AB0104211659548A728A849A47A0b07"
	block := uint64(30253895)
	BlockIndex := uint(229)
	TxIndex := uint(800)
	expectAmt0Wei, _ := decimal.NewFromString("99999999999999999999999972155")
	expectAmt1Wei, _ := decimal.NewFromString("0")
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

func TestMint_PancakeV3(t *testing.T) {
	// https://basescan.org/tx/0x0809776ecaeb651bd0f354e85814e7ed792709e5397ab777044a82a4a760d445#eventlog#296
	txHash := "0x0809776ecaeb651bd0f354e85814e7ed792709e5397ab777044a82a4a760d445"
	logIndex := 5
	pairAddress := "0x5845A51630AFab7C68556CF57a7b6827Bd94d434"
	token0Address := "0x858A6594f86fafb10dC0dEdC588C7CBb8E795129"
	block := uint64(30250940)
	BlockIndex := uint(99)
	TxIndex := uint(296)
	expectAmt0Wei, _ := decimal.NewFromString("97339241011744")
	expectAmt1Wei, _ := decimal.NewFromString("119999999999999")
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
