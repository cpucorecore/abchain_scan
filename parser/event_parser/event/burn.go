package event

import (
	"abchain_scan/repository/orm"
	"abchain_scan/types"
	"github.com/shopspring/decimal"
	"math/big"
)

type BurnEvent struct {
	*types.EventCommon
	Amount0Wei *big.Int
	Amount1Wei *big.Int
}

func (e *BurnEvent) CanGetTx() bool {
	return true
}

func (e *BurnEvent) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Event:         types.Remove,
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

var _ types.Event = (*BurnEvent)(nil)
