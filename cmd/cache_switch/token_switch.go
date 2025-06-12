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

type TokenSwitch struct {
	redisClient *redis.Client
	cache       cache.Cache
	workPool    *ants.Pool
	ctx         context.Context
}

func NewTokenSwitch(redisClient *redis.Client, cache cache.Cache, workPoolSize int) *TokenSwitch {
	workPool, err := ants.NewPool(workPoolSize)
	if err != nil {
		panic(err)
	}

	return &TokenSwitch{
		redisClient: redisClient,
		cache:       cache,
		workPool:    workPool,
		ctx:         context.Background(),
	}
}

func (s *TokenSwitch) getOldToken(key string) (*old_cache_types.Token, error) {
	value, err := s.redisClient.Get(s.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	token := &old_cache_types.Token{}
	if err = json.Unmarshal([]byte(value), token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *TokenSwitch) setNewToken(token *types.Token) {
	s.cache.SetToken(token)
}

func (s *TokenSwitch) switchToken(key string) {
	token, err := s.getOldToken(key)
	if err != nil {
		return
	}

	newToken := token.ToNewToken()
	s.setNewToken(newToken)
}

func (s *TokenSwitch) SwitchTokens() {
	log.Logger.Info("load token keys")
	keys, err := s.redisClient.Keys(s.ctx, "t:*").Result()
	if err != nil {
		return
	}
	log.Logger.Info("load token keys done", zap.Int("count", len(keys)))

	wg := &sync.WaitGroup{}
	for i, key := range keys {
		if i%1000 == 0 {
			log.Logger.Info("switch token", zap.Int("index", i))
		}
		wg.Add(1)
		_ = s.workPool.Submit(func() {
			defer wg.Done()
			s.switchToken(key)
		})
	}

	wg.Wait()
}
