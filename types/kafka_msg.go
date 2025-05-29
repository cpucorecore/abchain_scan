package types

import (
	"abchain_scan/repository/orm"
)

type BlockInfo struct {
	Height               uint64
	Timestamp            uint64
	NativeTokenPrice     string
	Txs                  []*orm.Tx
	NewTokens            []*orm.Token
	NewPairs             []*orm.Pair
	PoolUpdates          []*PoolUpdate
	PoolUpdateParameters []*PoolUpdateParameter
}

func (b *BlockInfo) ConvertABChainAddress() *BlockInfo {
	for _, tx := range b.Txs {
		tx.ConvertABChainAddress()
	}
	for _, token := range b.NewTokens {
		token.ConvertABChainAddress()
	}
	for _, pair := range b.NewPairs {
		pair.ConvertABChainAddress()
	}
	for _, poolUpdate := range b.PoolUpdates {
		poolUpdate.ConvertABChainAddress()
	}
	return b
}
