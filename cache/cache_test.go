package cache

import (
	"abchain_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

func TestCacheSetPrice(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewTwoTierCache(redisClient)

	blockNumber := big.NewInt(1)
	price, _ := decimal.NewFromString("33.33")
	cache.SetPrice(blockNumber, price)
}

func TestCacheGetPrice(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewTwoTierCache(redisClient)

	blockNumber := big.NewInt(1)
	price, _ := decimal.NewFromString("33.33")
	getPrice, b := cache.GetPrice(blockNumber)
	require.True(t, b)
	require.Equal(t, price, getPrice)
}

func TestCacheSetToken(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewTwoTierCache(redisClient)

	address := common.HexToAddress("0xe76004cffcab665c4692f663b8fb2a2f66adda9b")
	token := &types.Token{
		Address:     address,
		Creator:     address,
		Name:        "test",
		Symbol:      "test",
		Decimals:    18,
		TotalSupply: decimal.NewFromInt(1),
		BlockNumber: 1,
		BlockTime:   time.Unix(1000, 1),
		Program:     "test",
	}

	cache.SetToken(token)
}

func TestCacheGetToken(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := NewTwoTierCache(redisClient)

	address := common.HexToAddress("0xe76004cffcab665c4692f663b8fb2a2f66adda9b")
	expectToken := &types.Token{
		Address:     address,
		Creator:     address,
		Name:        "test",
		Symbol:      "test",
		Decimals:    18,
		TotalSupply: decimal.NewFromInt(1),
		BlockNumber: 1,
		BlockTime:   time.Unix(1000, 1),
		Program:     "test",
	}

	token, _ := cache.GetToken(address)
	require.True(t, token.Equal(expectToken))
}
