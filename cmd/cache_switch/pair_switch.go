package main

import (
	"abchain_scan/cache"
	"abchain_scan/cmd/cache_switch/old_cache_types"
	"abchain_scan/log"
	"abchain_scan/types"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"sync"
)

type PairSwitch struct {
	redisClient *redis.Client
	cache       cache.Cache
	workPool    *ants.Pool
	ctx         context.Context
}

func NewPairSwitch(redisClient *redis.Client, cache cache.Cache, workPoolSize int) *PairSwitch {
	workPool, err := ants.NewPool(workPoolSize)
	if err != nil {
		panic(err)
	}

	return &PairSwitch{
		redisClient: redisClient,
		cache:       cache,
		workPool:    workPool,
		ctx:         context.Background(),
	}
}

func (s *PairSwitch) getOldPair(key string) (*old_cache_types.Pair, error) {
	value, err := s.redisClient.Get(s.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	pair := &old_cache_types.Pair{}
	if err = json.Unmarshal([]byte(value), pair); err != nil {
		return nil, err
	}

	return pair, nil
}

func (s *PairSwitch) setNewPair(pair *types.Pair) {
	s.cache.SetPair(pair)
}

func (s *PairSwitch) switchPair(key string) {
	pair, err := s.getOldPair(key)
	if err != nil {
		return
	}

	newToken := pair.ToNewPair()
	s.setNewPair(newToken)
}

func (s *PairSwitch) SwitchPairs() {
	log.Logger.Info("load pair keys")
	keys, err := s.redisClient.Keys(s.ctx, "pr:*").Result()
	if err != nil {
		return
	}
	log.Logger.Info("load pair keys done", zap.Int("len", len(keys)))

	wg := &sync.WaitGroup{}
	for i, key := range keys {
		if i%1000 == 0 {
			log.Logger.Info("switch pair", zap.Int("i", i))
		}
		wg.Add(1)
		_ = s.workPool.Submit(func() {
			defer wg.Done()
			s.switchPair(key)
		})
	}

	wg.Wait()
}
