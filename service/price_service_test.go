package service

import (
	"abchain_scan/cache"
	"abchain_scan/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
)

func TestPriceService_GetBNBPrice(t *testing.T) {
	c := cache.NewMockCache()

	ethClient, err := ethclient.Dial(config.G.Chain.EndpointArchive)
	if err != nil {
		t.Fatal(err)
	}

	cc := NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())

	ps := NewPriceService(c, cc, ethClient, 0)
	price, err := ps.GetNativeTokenPrice(big.NewInt(67492426))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(price)
}
