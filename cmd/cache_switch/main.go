package main

import (
	"abchain_scan/cache"
	"abchain_scan/config"
	"flag"
	"github.com/go-redis/redis/v8"
)

func main() {
	config.LoadConfigFile("config.json")

	switchToken := false
	switchPair := false
	flag.BoolVar(&switchToken, "t", false, "switch token cache")
	flag.BoolVar(&switchPair, "p", false, "switch pair cache")
	flag.Parse()

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.G.Redis.Addr,
	})
	defer redisClient.Close()

	cache := cache.NewTwoTierCache(redisClient)

	if switchToken {
		tokenSwitch := NewTokenSwitch(redisClient, cache, 20)
		tokenSwitch.SwitchTokens()
	}

	if switchPair {
		pairSwitch := NewPairSwitch(redisClient, cache, 20)
		pairSwitch.SwitchPairs()
	}
}
