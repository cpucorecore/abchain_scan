package event

import (
	"abchain_scan/repository/orm"
	"abchain_scan/types"
	"github.com/shopspring/decimal"
	"math/big"
)

type MintEvent struct {
	*types.EventCommon
	Amount0Wei *big.Int
	Amount1Wei *big.Int
}

func (e *MintEvent) GetMintAmount() (decimal.Decimal, decimal.Decimal) {
	return ParseAmountsByPair(e.Amount0Wei, e.Amount1Wei, e.Pair)
}

func (e *MintEvent) CanGetTx() bool {
	return true
}

func (e *MintEvent) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Event:         types.Add,
		Maker:         e.Maker.String(),
		Token0Address: e.Pair.Token0Core.Address.String(),
		Token1Address: e.Pair.Token1Core.Address.String(),
		Block:         e.BlockNumber,
		BlockAt:       e.BlockTime,
		BlockIndex:    e.TxIndex,
		TxIndex:       e.LogIndex,
		PairAddress:   e.Pair.Address.String(),
		Program:       types.GetProtocolName(e.Pair.ProtocolId),
	}

	tx.Token0Amount, tx.Token1Amount = ParseAmountsByPair(e.Amount0Wei, e.Amount1Wei, e.Pair)
	tx.AmountUsd, tx.PriceUsd = CalcAmountAndPrice(bnbPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1Core.Address)
	return tx
}

func (e *MintEvent) IsMint() bool {
	return true
}

var _ types.Event = (*MintEvent)(nil)
