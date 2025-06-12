package service

import (
	"abchain_scan/abi/bep20"
	uniswapv2 "abchain_scan/abi/uniswap/v2"
	uniswapv3 "abchain_scan/abi/uniswap/v3"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type CallContractReq struct {
	BlockNumber *big.Int
	Address     *common.Address
	Data        []byte
}

func (r *CallContractReq) String() string {
	return "CallContractReq{" +
		"BlockNumber=" + r.BlockNumber.String() +
		", Address=" + r.Address.String() +
		", Data=" + hex.EncodeToString(r.Data) +
		"}"
}

func BuildCallContractReqDynamic(blockNumber *big.Int, address *common.Address, abi *abi.ABI, name string, args ...interface{}) *CallContractReq {
	data, err := abi.Pack(name, args...)
	if err != nil {
		panic(err)
	}

	return &CallContractReq{
		BlockNumber: blockNumber,
		Address:     address,
		Data:        data,
	}
}

var (
	AbiNames = []struct {
		Abi   *abi.ABI
		Names []string
	}{
		{
			Abi:   bep20.Abi,
			Names: []string{"name", "symbol", "decimals", "totalSupply"},
		},
		{
			Abi:   uniswapv2.PairAbi,
			Names: []string{"token0", "token1"},
		},
		{
			Abi:   uniswapv3.PoolAbi,
			Names: []string{"fee"},
		},
	}

	Name2Data map[string][]byte
)

func init() {
	Name2Data = make(map[string][]byte)
	for _, abiName := range AbiNames {
		for _, name := range abiName.Names {
			data, err := abiName.Abi.Pack(name)
			if err != nil {
				panic(err)
			}
			Name2Data[name] = data
		}
	}
}
