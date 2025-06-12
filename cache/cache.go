package cache

import (
	"abchain_scan/chain"
	"abchain_scan/log"
	"abchain_scan/types"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
	"time"
)

type PriceCache interface {
	SetPrice(blockNumber *big.Int, price decimal.Decimal)
	GetPrice(blockNumber *big.Int) (decimal.Decimal, bool)
}

type TokenCache interface {
	SetToken(token *types.Token)
	GetToken(address common.Address) (*types.Token, bool)
	DelToken(address common.Address)
}

type PairCache interface {
	SetPair(pair *types.Pair)
	GetPair(address common.Address) (*types.Pair, bool)
	PairExist(address common.Address) bool
	DelPair(address common.Address)
}

type BlockCache interface {
	SetFinishedBlock(blockNumber uint64)
	GetFinishedBlock() uint64
}

type Cache interface {
	PriceCache
	TokenCache
	PairCache
	BlockCache
}

type twoTierCache struct {
	ctx    context.Context
	memory *cache.Cache
	redis  *redis.Client
}

func NewTwoTierCache(redis *redis.Client) Cache {
	return &twoTierCache{
		ctx:    context.Background(),
		memory: cache.New(time.Hour*24, time.Hour),
		redis:  redis,
	}
}

func PriceCacheKey(blockNumber *big.Int) string {
	return fmt.Sprintf("%d:P:%s", chain.Id, blockNumber.String())
}

func TokenCacheKey(address common.Address) string {
	return fmt.Sprintf("%d:t:%s", chain.Id, address.Hex())
}

func PairCacheKey(address common.Address) string {
	return fmt.Sprintf("%d:p:%s", chain.Id, address.Hex())
}

var (
	fbKey = fmt.Sprintf("%d:fb", chain.Id)
)

func (c *twoTierCache) SetPrice(blockNumber *big.Int, price decimal.Decimal) {
	k := PriceCacheKey(blockNumber)
	c.memory.Set(k, price, cache.DefaultExpiration)
	err := c.redis.Set(c.ctx, k, price.String(), 0).Err()
	if err != nil {
		log.Logger.Error("save price failed", zap.Error(err))
	}
}

func (c *twoTierCache) GetPrice(blockNumber *big.Int) (decimal.Decimal, bool) {
	k := PriceCacheKey(blockNumber)
	price, ok := c.memory.Get(k)
	if ok {
		return price.(decimal.Decimal), true
	}

	v, err := c.redis.Get(c.ctx, k).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Logger.Error("get price failed", zap.Error(err))
		}
		return decimal.Zero, false
	}

	decimalPrice, err := decimal.NewFromString(v)
	if err != nil {
		return decimal.Decimal{}, false
	}
	c.memory.Set(k, decimalPrice, 0)
	return decimalPrice, true
}

func (c *twoTierCache) SetToken(token *types.Token) {
	token.Timestamp = time.Now()
	k := TokenCacheKey(token.Address)
	c.memory.Set(k, token, cache.DefaultExpiration)
	err := c.redis.Set(c.ctx, k, token, 0).Err()
	if err != nil {
		log.Logger.Error("save token failed", zap.Error(err))
	}
}

func (c *twoTierCache) GetToken(address common.Address) (*types.Token, bool) {
	k := TokenCacheKey(address)
	tokenCache, ok := c.memory.Get(k)
	if ok {
		return tokenCache.(*types.Token), true
	}

	v := &types.Token{}
	err := c.redis.Get(c.ctx, k).Scan(v)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Logger.Error("redis get err", zap.Error(err))
		}
		return nil, false
	}
	return v, true
}

func (c *twoTierCache) DelToken(address common.Address) {
	k := TokenCacheKey(address)
	c.memory.Delete(k)
	err := c.redis.Del(c.ctx, k).Err()
	if err != nil {
		log.Logger.Error("redis del err", zap.Error(err))
	}
}

func (c *twoTierCache) SetPair(pair *types.Pair) {
	pair.Timestamp = time.Now()
	k := PairCacheKey(pair.Address)
	c.memory.Set(k, pair, cache.DefaultExpiration)
	err := c.redis.Set(c.ctx, k, pair, 0).Err()
	if err != nil {
		log.Logger.Error("save pair failed", zap.Error(err))
	}
}

func (c *twoTierCache) GetPair(address common.Address) (*types.Pair, bool) {
	k := PairCacheKey(address)
	pair, ok := c.memory.Get(k)
	if ok {
		return pair.(*types.Pair), true
	}

	v := &types.Pair{}
	err := c.redis.Get(c.ctx, k).Scan(v)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Logger.Error("redis get err", zap.Error(err))
		}
		return nil, false
	}

	c.memory.Set(k, v, 0)
	return v, true
}

func (c *twoTierCache) PairExist(address common.Address) bool {
	_, exist := c.GetPair(address)
	return exist
}

func (c *twoTierCache) DelPair(address common.Address) {
	k := PairCacheKey(address)
	c.memory.Delete(k)
	err := c.redis.Del(c.ctx, k).Err()
	if err != nil {
		log.Logger.Error("redis del err", zap.Error(err))
	}
}

func (c *twoTierCache) SetFinishedBlock(blockNumber uint64) {
	c.redis.Set(c.ctx, fbKey, blockNumber, 0)
}

func (c *twoTierCache) GetFinishedBlock() uint64 {
	v, err := c.redis.Get(c.ctx, fbKey).Uint64()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Logger.Error("redis get err", zap.Error(err))
		}
		return 0
	}
	return v
}
