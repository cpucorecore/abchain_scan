package service

import (
	"abchain_scan/cache"
	"abchain_scan/log"
	"abchain_scan/metrics"
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
	"time"
)

type PriceService interface {
	Start(startBlockNumber uint64)
	GetNativeTokenPrice(blockNumber *big.Int) (decimal.Decimal, error)
}

type priceService struct {
	cache          cache.Cache
	contractCaller *ContractCaller
	workPoolSize   int
	workPool       *ants.Pool
	ethClient      *ethclient.Client
}

func NewPriceService(
	cache cache.Cache,
	contractCaller *ContractCaller,
	ethClient *ethclient.Client,
	poolSize int,
) PriceService {
	var workPool *ants.Pool
	var err error
	if poolSize > 0 {
		workPool, err = ants.NewPool(poolSize)
		if err != nil {
			log.Logger.Fatal("ants pool(BlockGetter) init err", zap.Error(err))
		}
	}

	return &priceService{
		cache:          cache,
		contractCaller: contractCaller,
		workPoolSize:   poolSize,
		workPool:       workPool,
		ethClient:      ethClient,
	}
}

func (ps *priceService) Start(startBlockNumber uint64) {
	if ps.workPoolSize <= 0 {
		return
	}

	go func() {
		for {
			headerBlockNumber, err := ps.ethClient.BlockNumber(context.Background())
			if err != nil {
				log.Logger.Error("ethClient.HeightBigInt", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}

			for startBlockNumber <= headerBlockNumber {
				ps.workPool.Submit(func() {
					ps.GetNativeTokenPrice(big.NewInt(int64(startBlockNumber)))
					startBlockNumber++
				})
			}
		}
	}()
}

func (ps *priceService) GetNativeTokenPrice(blockNumber *big.Int) (decimal.Decimal, error) {
	// mock
	return decimal.NewFromBigInt(blockNumber, 1), nil

	cachePrice, ok := ps.cache.GetPrice(blockNumber)
	if ok {
		return cachePrice, nil
	}

	return ps.getNativeTokenPrice(blockNumber)
}

func (ps *priceService) getNativeTokenPrice(blockNumber *big.Int) (decimal.Decimal, error) {
	now := time.Now()
	reserve0, reserve1, err := ps.contractCaller.GetReservesByBlockNumber(blockNumber)
	if err != nil {
		log.Logger.Error("GetReservesByBlockNumber err", zap.Error(err), zap.Uint64("blockNumber", blockNumber.Uint64()))
		return decimal.Zero, err
	}
	metrics.CallContractArchiveDurationMs.Observe(float64(time.Since(now).Milliseconds()))

	// TODO get pair info
	USDCAmountDivWETHAmount := decimal.NewFromBigInt(reserve1, -6).Div(decimal.NewFromBigInt(reserve0, -18))
	ps.cache.SetPrice(blockNumber, USDCAmountDivWETHAmount)

	return USDCAmountDivWETHAmount, nil
}
