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

func TestBurn_Aerodrome(t *testing.T) {
	// https://basescan.org/tx/0xeb50cf26b45a8ca72d9343dab433f5eeecaeb3eabc2716a7dd80ddef966948b1#eventlog#213
	tc := service.GetTestContext()
	ethLog := tc.GetEthLog("0xeb50cf26b45a8ca72d9343dab433f5eeecaeb3eabc2716a7dd80ddef966948b1", 6)

	event, pErr := Topic2EventParser[ethLog.Topics[0]].Parse(ethLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPair(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(service.MockNativeTokenPrice)
	expectAmt0, _ := decimal.NewFromString("4047.640731408680145311")
	expectAmt1, _ := decimal.NewFromString("9.88229999999999995")
	expectTx := &orm.Tx{
		TxHash:        "0xeb50cf26b45a8ca72d9343dab433f5eeecaeb3eabc2716a7dd80ddef966948b1",
		Event:         types.Remove,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Token0Address: "0x1E50309675d5C41D38Ba14133B4DB5b9f44FfBCd",
		Token1Address: types.WETH,
		Block:         30233109,
		BlockIndex:    66,
		TxIndex:       213,
		PairAddress:   "0xC09F68906B1DC60F1BA5771Ec6625cA947031Aaf",
		Program:       types.ProtocolNameAerodrome,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}

func TestBurn_UniswapV2(t *testing.T) {
	// https://basescan.org/tx/0xd9f86ab4778970d495ab09e3aafff24bc4ecf55f2e99d5e79cc8dc19e36fb401#eventlog#1051
	txHash := "0xd9f86ab4778970d495ab09e3aafff24bc4ecf55f2e99d5e79cc8dc19e36fb401"
	logIndex := 5
	pairAddress := "0x918B7b442862C9A39814CFdAc197E43490f73A01"
	token0Address := "0xea778385f93e3184E1F33270C0BA10eD6D805834"
	block := uint64(30294196)
	BlockIndex := uint(285)
	TxIndex := uint(1051)
	expectAmt0Wei, _ := decimal.NewFromString("228801199741494110064733")
	expectAmt1Wei, _ := decimal.NewFromString("5524217866093027057")
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
		Event:         types.Remove,
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

func TestBurn_PancakeV2(t *testing.T) {
	// TODO
}
