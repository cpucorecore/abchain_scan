package service

import (
	uniswapv2 "abchain_scan/abi/uniswap/v2"
	uniswapv3 "abchain_scan/abi/uniswap/v3"
	"abchain_scan/cache"
	"abchain_scan/log"
	"abchain_scan/metrics"
	"abchain_scan/types"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"math/big"
	"sync"
	"time"
)

type PairService interface {
	SetPair(pair *types.Pair)
	GetPairTokens(pair *types.Pair) *types.PairWrap
	GetPair(pairAddress common.Address, possibleProtocolIds []int) *types.PairWrap
}

type pairService struct {
	ctx            context.Context
	cache          cache.Cache
	contractCaller *ContractCaller
	group          singleflight.Group
}

func NewPairService(
	cache cache.Cache,
	contractCaller *ContractCaller,
) PairService {
	return &pairService{
		ctx:            context.Background(),
		cache:          cache,
		contractCaller: contractCaller,
	}
}

func (s *pairService) SetPair(pair *types.Pair) {
	s.cache.SetPair(pair)
}

func (s *pairService) doGetToken(tokenAddress common.Address) (*types.Token, error) {
	token := &types.Token{
		Address: tokenAddress,
	}

	var (
		nameRes struct {
			name string
			err  error
		}
		symbolRes struct {
			symbol string
			err    error
		}
		decimalsRes struct {
			decimals int
			err      error
		}
		supplyRes struct {
			supply *big.Int
			err    error
		}
	)

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		nameRes.name, nameRes.err = s.contractCaller.CallName(&tokenAddress)
	}()

	go func() {
		defer wg.Done()
		symbolRes.symbol, symbolRes.err = s.contractCaller.CallSymbol(&tokenAddress)
	}()

	go func() {
		defer wg.Done()
		decimalsRes.decimals, decimalsRes.err = s.contractCaller.CallDecimals(&tokenAddress)
	}()

	go func() {
		defer wg.Done()
		supplyRes.supply, supplyRes.err = s.contractCaller.CallTotalSupply(&tokenAddress)
	}()
	wg.Wait()

	if nameRes.err == nil {
		token.Name = nameRes.name
	}

	if symbolRes.err == nil {
		token.Symbol = symbolRes.symbol
	}

	if decimalsRes.err != nil {
		token.Filtered = true
		return token, decimalsRes.err
	}
	token.Decimals = int8(decimalsRes.decimals)

	if supplyRes.err == nil {
		token.TotalSupply = decimal.NewFromBigInt(supplyRes.supply, -int32(token.Decimals))
	}

	return token, nil
}

func (s *pairService) getToken(tokenAddress common.Address) (*types.Token, error, bool) {
	cacheToken, ok := s.cache.GetToken(tokenAddress)
	if ok {
		return cacheToken, nil, true
	}

	now := time.Now()
	doResult, err, _ := s.group.Do(tokenAddress.String(), func() (interface{}, error) {
		token, err := s.doGetToken(tokenAddress)
		s.cache.SetToken(token)
		return token, err
	})
	if err != nil {
		return nil, err, false
	}

	metrics.GetTokenDurationMs.Observe(float64(time.Since(now).Milliseconds()))
	return doResult.(*types.Token), nil, false
}

func (s *pairService) getPairTokens(pair *types.Pair) *types.PairWrap {
	pairWrap := &types.PairWrap{
		Pair:    pair,
		NewPair: !s.cache.PairExist(pair.Address),
	}

	var (
		wg                               sync.WaitGroup
		token0                           *types.Token
		token1                           *types.Token
		token0Err, token1Err             error
		token0FromCache, token1FromCache bool
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		t0, err, fromCache := s.getToken(pair.Token0Core.Address)
		token0, token0Err, token0FromCache = t0, err, fromCache
	}()

	go func() {
		defer wg.Done()
		t1, err, fromCache := s.getToken(pair.Token1Core.Address)
		token1, token1Err, token1FromCache = t1, err, fromCache
	}()
	wg.Wait()

	if token0Err != nil {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken0
		return pairWrap
	}

	if token1Err != nil {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken1
		return pairWrap
	}

	pair.Token0 = token0
	pair.Token1 = token1

	pair.Token0Core.Symbol = token0.Symbol
	pair.Token0Core.Decimals = token0.Decimals
	pair.Token1Core.Symbol = token1.Symbol
	pair.Token1Core.Decimals = token1.Decimals

	tokensReversed := pair.OrderToken0Token1()
	if tokensReversed {
		pairWrap.NewToken0 = !token1FromCache
		pairWrap.NewToken1 = !token0FromCache
	} else {
		pairWrap.NewToken0 = !token0FromCache
		pairWrap.NewToken1 = !token1FromCache
	}

	return pairWrap
}

func (s *pairService) GetPairTokens(pair *types.Pair) *types.PairWrap {
	doResult, _, _ := s.group.Do(pair.Address.String(), func() (interface{}, error) {
		pairWrap := s.getPairTokens(pair)
		s.SetPair(pair)
		return pairWrap, nil
	})

	return doResult.(*types.PairWrap)
}

func (s *pairService) getPair(pairAddress common.Address, possibleProtocolIds []int) *types.PairWrap {
	doResult, _, _ := s.group.Do(pairAddress.String()+"gp", func() (interface{}, error) {
		pair := s.doGetPair(pairAddress)
		if pair.Filtered {
			s.SetPair(pair)
			return &types.PairWrap{
				Pair:      pair,
				NewPair:   false,
				NewToken0: false,
				NewToken1: false,
			}, nil
		}

		if !s.verifyPair(pair, possibleProtocolIds) {
			s.SetPair(pair)
			return &types.PairWrap{
				Pair:      pair,
				NewPair:   false,
				NewToken0: false,
				NewToken1: false,
			}, nil
		}

		return s.GetPairTokens(pair), nil
	})

	return doResult.(*types.PairWrap)
}

func (s *pairService) GetPair(pairAddress common.Address, possibleProtocolIds []int) *types.PairWrap {
	cachePair, ok := s.cache.GetPair(pairAddress)
	if ok {
		return &types.PairWrap{
			Pair: cachePair,
		}
	}

	return s.getPair(pairAddress, possibleProtocolIds)
}

func (s *pairService) doGetPair(pairAddress common.Address) *types.Pair {
	pair := &types.Pair{
		Address: pairAddress,
	}

	var (
		token0Res struct {
			address common.Address
			err     error
		}
		token1Res struct {
			address common.Address
			err     error
		}
	)

	var wg sync.WaitGroup
	wg.Add(2)
	now := time.Now()
	go func() {
		defer wg.Done()
		token0Res.address, token0Res.err = s.contractCaller.CallToken0(&pairAddress)
	}()

	go func() {
		defer wg.Done()
		token1Res.address, token1Res.err = s.contractCaller.CallToken1(&pairAddress)
	}()
	wg.Wait()

	if token0Res.err != nil {
		log.Logger.Info("Err: CallToken0 err, this pair will filtered",
			zap.Error(token0Res.err),
			zap.String("pair address", pairAddress.String()),
		)
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken0
		return pair
	}
	pair.Token0Core = &types.TokenCore{
		Address: token0Res.address,
	}

	if token1Res.err != nil {
		log.Logger.Info("Err: CallToken1 err, this pair will filtered",
			zap.Error(token1Res.err),
			zap.String("pair address", pairAddress.String()),
		)
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken1
		return pair
	}
	pair.Token1Core = &types.TokenCore{
		Address: token1Res.address,
	}

	metrics.GetPairDurationMs.Observe(float64(time.Since(now).Milliseconds()))
	pair.FilterByToken0AndToken1()
	return pair
}

func (s *pairService) verifyPairV2(pairFactoryAddress common.Address, pair *types.Pair) bool {
	pairAddressQueried, err := s.contractCaller.CallGetPair(&pairFactoryAddress, &pair.Token0Core.Address, &pair.Token1Core.Address)
	if err != nil {
		return false
	}
	return types.IsSameAddress(pairAddressQueried, pair.Address)
}

func (s *pairService) verifyPairV3(pairFactoryAddress common.Address, pair *types.Pair) bool {
	fee, callFeeErr := s.contractCaller.CallFee(&pair.Address)
	if callFeeErr != nil {
		return false
	}

	pairAddressQueried, err := s.contractCaller.CallGetPool(&pairFactoryAddress, &pair.Token0Core.Address, &pair.Token1Core.Address, fee)
	if err != nil {
		return false
	}

	return types.IsSameAddress(pairAddressQueried, pair.Address)
}

func (s *pairService) verifyPair(pair *types.Pair, possibleProtocolIds []int) bool {
	now := time.Now()
	defer func() {
		duration := float64(time.Since(now).Milliseconds())
		metrics.VerifyPairDurationMs.Observe(duration)
	}()

	for _, protocolId := range possibleProtocolIds {
		switch protocolId {
		case types.ProtocolIdNewSwap:
			if s.verifyPairV2(uniswapv2.FactoryAddress, pair) {
				pair.ProtocolId = protocolId
				metrics.VerifyPairTotal.WithLabelValues("success").Inc()
				metrics.VerifyPairOkByProtocol.WithLabelValues("uniswap_v2").Inc()
				return true
			}

		case types.ProtocolIdUniswapV3:
			if s.verifyPairV3(uniswapv3.FactoryAddress, pair) {
				pair.ProtocolId = protocolId
				metrics.VerifyPairTotal.WithLabelValues("success").Inc()
				metrics.VerifyPairOkByProtocol.WithLabelValues("uniswap_v3").Inc()
				return true
			}
		}
	}

	pair.Filtered = true
	pair.FilterCode = types.FilterCodeVerifyFailed
	metrics.VerifyPairTotal.WithLabelValues("failed").Inc()

	return false
}
