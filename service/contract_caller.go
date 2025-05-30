package service

import (
	"abchain_scan/abi/erc20"
	uniswapv2 "abchain_scan/abi/uniswap/v2"
	"abchain_scan/config"
	"abchain_scan/metrics"
	"abchain_scan/types"
	"context"
	"errors"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
	"time"
)

var (
	ErrOutputEmpty       = errors.New("output is empty")
	ErrWrongOutputLength = errors.New("wrong output length")
	ErrReserve0NotBigInt = errors.New("reverse0 is not *big.Int")
	ErrReserve1NotBigInt = errors.New("reverse1 is not *big.Int")
)

type ContractCaller struct {
	ctx         context.Context
	ethClient   *ethclient.Client
	retryParams *config.RetryParams
}

func NewContractCaller(ethClient *ethclient.Client, retryParams *config.RetryParams) *ContractCaller {
	return &ContractCaller{
		ctx:         context.Background(),
		ethClient:   ethClient,
		retryParams: retryParams,
	}
}

func IsRetryableErr(err error) bool {
	errMsg := err.Error()
	if strings.Contains(errMsg, "execution reverted") ||
		strings.Contains(errMsg, "out of gas") ||
		strings.Contains(errMsg, "abi: cannot marshal in to go slice") {
		return false
	}
	return true
}

func (c *ContractCaller) callContract(req *CallContractReq) ([]byte, error) {
	now := time.Now()
	bytes, err := c.ethClient.CallContract(
		c.ctx,
		ethereum.CallMsg{
			To:   req.Address,
			Data: req.Data,
		},
		req.BlockNumber,
	)

	if err != nil {
		if IsRetryableErr(err) {
			metrics.CallContractErrors.WithLabelValues("true").Inc()
			//log.Logger.Info("Err: call contract encounter retryable err", zap.Error(err), zap.Any("req", req))
			return nil, err
		}

		metrics.CallContractErrors.WithLabelValues("false").Inc()
		//log.Logger.Info("Err: call contract encounter no retryable err", zap.Error(err), zap.Any("req", req))
		return nil, nil
	}

	metrics.CallContractDurationMs.Observe(float64(time.Since(now).Milliseconds()))

	return bytes, nil
}

func (c *ContractCaller) CallContract(req *CallContractReq) ([]byte, error) {
	ctxWithTimeout, _ := context.WithTimeout(c.ctx, c.retryParams.Timeout)
	return retry.DoWithData(func() ([]byte, error) {
		return c.callContract(req)
	}, c.retryParams.Attempts, c.retryParams.Delay, retry.Context(ctxWithTimeout))
}

func (c *ContractCaller) getString(address *common.Address, method string) (string, error) {
	req := &CallContractReq{
		Address: address,
		Data:    Name2Data[method],
	}

	bytes, err := c.CallContract(req)
	if err != nil {
		return "", err
	}

	if len(bytes) == 0 {
		return "", ErrOutputEmpty
	}

	var stringValue string
	err = erc20.Abi.Unpack(&stringValue, method, bytes)
	if err != nil {
		return "", err
	}

	return stringValue, nil
}

func (c *ContractCaller) CallName(address *common.Address) (string, error) {
	return c.getString(address, "name")
}

func (c *ContractCaller) CallSymbol(address *common.Address) (string, error) {
	return c.getString(address, "symbol")
}

func (c *ContractCaller) CallDecimals(address *common.Address) (int, error) {
	method := "decimals"
	req := &CallContractReq{
		Address: address,
		Data:    Name2Data[method],
	}

	bytes, err := c.CallContract(req)
	if err != nil {
		return 0, err
	}

	if len(bytes) == 0 {
		return 0, ErrOutputEmpty
	}

	var value uint8
	err = erc20.Abi.Unpack(&value, method, bytes)
	if err != nil {
		return 0, err
	}

	return int(value), nil
}

func (c *ContractCaller) CallTotalSupply(address *common.Address) (*big.Int, error) {
	method := "totalSupply"
	req := &CallContractReq{
		Address: address,
		Data:    Name2Data[method],
	}

	bytes, err := c.CallContract(req)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, ErrOutputEmpty
	}

	var value *big.Int
	err = erc20.Abi.Unpack(&value, method, bytes)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (c *ContractCaller) getAddress(address *common.Address, method string) (*common.Address, error) {
	req := &CallContractReq{
		Address: address,
		Data:    Name2Data[method],
	}

	bytes, err := c.CallContract(req)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, ErrOutputEmpty
	}

	var value common.Address
	err = uniswapv2.PairAbi.Unpack(&value, method, bytes)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (c *ContractCaller) CallToken0(address *common.Address) (*common.Address, error) {
	return c.getAddress(address, "token0")
}

func (c *ContractCaller) CallToken1(address *common.Address) (*common.Address, error) {
	return c.getAddress(address, "token1")
}

/*
CallGetPair
for uniswap/pancake v2
*/
func (c *ContractCaller) CallGetPair(factoryAddress, token0Address, token1Address *common.Address) (common.Address, error) {
	req := BuildCallContractReqDynamic(nil, factoryAddress, uniswapv2.FactoryAbi, "getPair", token0Address, token1Address)

	bytes, err := c.CallContract(req)
	if err != nil {
		return types.ZeroAddress, err
	}

	if len(bytes) == 0 {
		return types.ZeroAddress, ErrOutputEmpty
	}

	var value common.Address
	err = uniswapv2.FactoryAbi.Unpack(&value, "getPair", bytes)
	if err != nil {
		return types.ZeroAddress, err
	}

	return value, nil
}

/*
callGetReserves
for uniswap/pancake v2
*/
func (c *ContractCaller) callGetReserves(blockNumber *big.Int) (map[string]interface{}, error) {
	req := BuildCallContractReqDynamic(blockNumber, &types.WETHUSDCPairAddressUniswapV2, uniswapv2.PairAbi, "getReserves")

	bytes, err := c.CallContract(req)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, ErrOutputEmpty
	}

	values, unpackErr := UniswapV2PairUnpacker.Unpack("getReserves", bytes, 3)
	if unpackErr != nil {
		return nil, unpackErr
	}

	if len(values) != 3 {
		return nil, ErrWrongOutputLength
	}

	return values, nil
}

func (c *ContractCaller) GetReservesByBlockNumber(blockNumber *big.Int) (*big.Int, *big.Int, error) {
	values, err := c.callGetReserves(blockNumber)
	if err != nil {
		return nil, nil, err
	}

	reserve0, ok0 := values["_reserve0"].(*big.Int)
	if !ok0 {
		return nil, nil, ErrReserve0NotBigInt
	}

	reserve1, ok1 := values["_reserve1"].(*big.Int)
	if !ok1 {
		return nil, nil, ErrReserve1NotBigInt
	}

	return reserve0, reserve1, nil
}
