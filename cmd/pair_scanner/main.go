package main

import (
	uniswapv2 "abchain_scan/abi/uniswap/v2"
	"abchain_scan/cache"
	"abchain_scan/config"
	"abchain_scan/log"
	"abchain_scan/service"
	"abchain_scan/types"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"math/big"
	"sync"
)

func main() {
	config.LoadConfigFile("config.json")
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.G.Redis.Addr,
		Username: config.G.Redis.Username,
		Password: config.G.Redis.Password,
	})
	c := cache.NewTwoTierCache(rdb)
	ethClient, err := ethclient.Dial(config.G.Chain.WsEndpoint)
	if err != nil {
		panic("Failed to connect to the chain: " + err.Error())
	}
	cc := service.NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())
	ps := service.NewPairService(c, cc)
	ctx := context.Background()
	data, err := uniswapv2.FactoryAbi.Pack("allPairsLength")
	if err != nil {
		panic("Failed to pack allPairsLength: " + err.Error())
	}
	cm := ethereum.CallMsg{
		To:   &uniswapv2.FactoryAddress,
		Data: data,
	}
	data, err = ethClient.CallContract(ctx, cm, nil)
	if err != nil {
		panic("Failed to call contract: " + err.Error())
	}
	length := new(big.Int)
	err = uniswapv2.FactoryAbi.Unpack(&length, "allPairsLength", data)
	if err != nil {
		panic("Failed to unpack allPairsLength: " + err.Error())
	}
	log.Logger.Info("Total pairs length", zap.Any("length", length))

	workPool, err := ants.NewPool(1) // use 1 to query contract in order, or you can use a larger number for parallel queries
	if err != nil {
		log.Logger.Fatal("Failed to create ants pool", zap.Error(err))
	}
	defer workPool.Release()

	wg := &sync.WaitGroup{}
	for i := 0; i < int(length.Int64()); i++ {
		wg.Add(1)
		workPool.Submit(func() {
			defer wg.Done()
			data, err = uniswapv2.FactoryAbi.Pack("allPairs", big.NewInt(int64(i)))
			if err != nil {
				log.Logger.Error("Failed to pack allPairs", zap.Error(err))
				return
			}
			cm.Data = data
			data, err = ethClient.CallContract(ctx, cm, nil)
			if err != nil {
				log.Logger.Error("Failed to call contract allPairs", zap.Error(err))
				return
			}
			var pa common.Address
			err = uniswapv2.FactoryAbi.Unpack(&pa, "allPairs", data)
			if err != nil {
				log.Logger.Error("Failed to unpack allPairs", zap.Error(err))
				return
			}

			pairWrap := ps.GetPair(pa, []int{types.ProtocolIdUniswapV2})
			pairInfo, _ := pairWrap.Pair.MarshalBinary()
			log.Logger.Info("Total pairs", zap.String("pair", string(pairInfo)))
		})
	}

	wg.Wait()
}
