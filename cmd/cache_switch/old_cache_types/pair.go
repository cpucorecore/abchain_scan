package old_cache_types

import (
	"abchain_scan/log"
	"abchain_scan/types"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
)

type TokenInfo struct {
	Address     common.Address
	Name        string
	Symbol      string
	Decimals    int16
	TotalSupply decimal.Decimal
}

type Pair struct {
	Address    common.Address
	Token0     *TokenInfo
	Token1     *TokenInfo
	Block      uint64
	BlockAt    time.Time
	TxIndex    uint
	TxHash     common.Hash
	From       common.Address
	ProtocolId int
	Filtered   bool
	FilterCode int
}

func (p *Pair) MarshalBinary() ([]byte, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (p *Pair) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &p)
}

func (p *Pair) ToNewPair() *types.Pair {
	pair := &types.Pair{
		Address:        p.Address,
		TokensReversed: false,
		Block:          p.Block,
		BlockAt:        p.BlockAt,
		ProtocolId:     p.ProtocolId,
		Filtered:       p.Filtered,
		FilterCode:     p.FilterCode,
	}

	if p.Token0 != nil {
		pair.Token0Core = &types.TokenCore{
			Address:  p.Token0.Address,
			Symbol:   p.Token0.Symbol,
			Decimals: int8(p.Token0.Decimals),
		}
	} else {
		log.Logger.Info("pair token0 is nil", zap.String("address", p.Address.String()))
	}

	if p.Token1 != nil {
		pair.Token1Core = &types.TokenCore{
			Address:  p.Token1.Address,
			Symbol:   p.Token1.Symbol,
			Decimals: int8(p.Token1.Decimals),
		}
	} else {
		log.Logger.Info("pair token1 is nil", zap.String("address", p.Address.String()))
	}

	swaped := pair.OrderToken0Token1()
	if swaped {
		log.Logger.Info("pair swaped", zap.String("address", p.Address.String()))
	}

	return pair
}
