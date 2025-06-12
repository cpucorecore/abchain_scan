package event

import (
	"abchain_scan/repository/orm"
	"abchain_scan/types"
	"github.com/shopspring/decimal"
	"math/big"
)

type SwapEventV3 struct {
	*types.EventCommon
	Amount0Wei *big.Int
	Amount1Wei *big.Int
}

func (e *SwapEventV3) CanGetTx() bool {
	return true
}

func (e *SwapEventV3) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
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
	if tx.Token0Amount.IsNegative() {
		tx.Event = types.Buy
		tx.Token0Amount = tx.Token0Amount.Neg()
	} else if tx.Token1Amount.IsNegative() {
		tx.Event = types.Sell
		tx.Token1Amount = tx.Token1Amount.Neg()
	} else {
	}

	tx.AmountUsd, tx.PriceUsd = CalcAmountAndPrice(bnbPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1Core.Address)
	return tx
}

func (e *SwapEventV3) CanGetPoolUpdateParameter() bool {
	return true
}

func (e *SwapEventV3) GetPoolUpdateParameter() *types.PoolUpdateParameter {
	return &types.PoolUpdateParameter{
		BlockNumber:   e.BlockNumber,
		PairAddress:   e.Pair.Address,
		Token0Address: e.Pair.Token0Core.Address,
		Token1Address: e.Pair.Token1Core.Address,
	}
}

var _ types.Event = (*SwapEventV3)(nil)
