package service

import (
	"github.com/shopspring/decimal"
	"math/big"
)

var mockPrice = decimal.NewFromFloat(10)

type priceServiceMock struct {
}

func NewPriceServiceMock() PriceService {
	return &priceServiceMock{}
}

func (psm *priceServiceMock) Start(startBlockNumber uint64) {
}

func (psm *priceServiceMock) GetNativeTokenPrice(blockNumber *big.Int) (decimal.Decimal, error) {
	return mockPrice, nil
}

var _ PriceService = (*priceServiceMock)(nil)
