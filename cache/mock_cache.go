package cache

import (
	"abchain_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"math/big"
)

type MockCache struct {
	memory *cache.Cache
}

func NewMockCache() Cache {
	return &MockCache{
		memory: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

func (c *MockCache) DelToken(address common.Address) {
	c.memory.Delete(address.String())
}

func (c *MockCache) DelPair(address common.Address) {
	c.memory.Delete(address.String())
}

func (c *MockCache) SetPrice(blockNumber *big.Int, price decimal.Decimal) {
	c.memory.Set(blockNumber.String(), price, 0)
}

func (c *MockCache) GetPrice(blockNumber *big.Int) (decimal.Decimal, bool) {
	if price, found := c.memory.Get(blockNumber.String()); found {
		return price.(decimal.Decimal), true
	}
	return decimal.Decimal{}, false
}

func (c *MockCache) SetToken(token *types.Token) {
	c.memory.Set(token.Address.String(), token, 0)
}

func (c *MockCache) GetToken(address common.Address) (*types.Token, bool) {
	if token, found := c.memory.Get(address.String()); found {
		return token.(*types.Token), true
	}
	return nil, false
}

func (c *MockCache) SetPair(pair *types.Pair) {
	c.memory.Set(pair.Address.String(), pair, 0)
}

func (c *MockCache) GetPair(address common.Address) (*types.Pair, bool) {
	if pair, found := c.memory.Get(address.String()); found {
		return pair.(*types.Pair), true
	}
	return nil, false
}

func (c *MockCache) PairExist(address common.Address) bool {
	if _, found := c.memory.Get(address.String()); found {
		return true
	}
	return false
}

func (c *MockCache) SetFinishedBlock(blockNumber uint64) {
}

func (c *MockCache) GetFinishedBlock() uint64 {
	return 0
}

var _ Cache = &MockCache{}
