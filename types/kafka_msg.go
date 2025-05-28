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

type BlockInfoOld struct {
	BlockNumber            uint64
	BlockAt                uint64
	BnbPrice               string
	Txs                    []*orm.Tx
	PoolUpdatesV2          []*PoolUpdate
	PoolUpdateParametersV3 []*PoolUpdateParameter
}
