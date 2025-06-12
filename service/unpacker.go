package service

import (
	"abchain_scan/abi/bep20"
	"abchain_scan/abi/ds_token"
	uniswapv2 "abchain_scan/abi/uniswap/v2"
	uniswapv3 "abchain_scan/abi/uniswap/v3"
	"abchain_scan/types"
	"bytes"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
	"unicode/utf8"
)

type Unpacker interface {
	Unpack(method string, data []byte, length int) (values []interface{}, err error)
}

type unpacker struct {
	abis []*abi.ABI
}

func NewUnpacker(abis []*abi.ABI) Unpacker {
	return &unpacker{
		abis: abis,
	}
}

var (
	UnpackErr           = errors.New("unpack error")
	ErrWrongString      = errors.New("wrong string")
	ErrWrongIntType     = errors.New("wrong int type")
	ErrWrongBigIntType  = errors.New("wrong big int type")
	ErrWrongAddressType = errors.New("wrong address type")
	ErrWrongBoolType    = errors.New("wrong bool type")
)

func (u *unpacker) Unpack(method string, data []byte, length int) ([]interface{}, error) {
	for _, abi_ := range u.abis {
		values, err := abi_.Unpack(method, data)
		if err == nil && len(values) == length {
			return values, nil
		}
	}
	return nil, UnpackErr
}

var (
	TokenUnpacker = NewUnpacker([]*abi.ABI{
		bep20.Abi,
		ds_token.Abi,
	})

	UniswapV2PairUnpacker = NewUnpacker([]*abi.ABI{
		uniswapv2.PairAbi,
	})

	UniswapV3PoolUnpacker = NewUnpacker([]*abi.ABI{
		uniswapv3.PoolAbi,
	})

	UniswapV2FactoryUnpacker = NewUnpacker([]*abi.ABI{
		uniswapv2.FactoryAbi,
	})

	UniswapV3FactoryUnpacker = NewUnpacker([]*abi.ABI{
		uniswapv3.FactoryAbi,
	})

	Name2Unpacker = map[string]Unpacker{
		"name":        TokenUnpacker,
		"symbol":      TokenUnpacker,
		"decimals":    TokenUnpacker,
		"totalSupply": TokenUnpacker,
		"token0":      UniswapV2PairUnpacker,
		"token1":      UniswapV2PairUnpacker,
		"getReserves": UniswapV2PairUnpacker,
		"fee":         UniswapV3PoolUnpacker,
	}
)

func sanitizeUTF8(s string) string {
	if !utf8.ValidString(s) {
		return strings.ToValidUTF8(s, "?")
	}
	return s
}

func ParseString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case [32]byte:
		return sanitizeUTF8(string(bytes.ReplaceAll(v[:], []byte{0}, []byte{}))), nil
	default:
		return "", ErrWrongString
	}
}

func ParseInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case uint8:
		return int(v), nil
	case *big.Int:
		return int(v.Int64()), nil
	default:
		return 0, ErrWrongIntType
	}
}

func ParseBigInt(value interface{}) (*big.Int, error) {
	if bigIntValue, ok := value.(*big.Int); ok {
		return bigIntValue, nil
	}
	return nil, ErrWrongBigIntType
}

func ParseAddress(value interface{}) (common.Address, error) {
	if address, ok := value.(common.Address); ok {
		return address, nil
	} else {
		return types.ZeroAddress, ErrWrongAddressType
	}
}

func ParseBool(value interface{}) (bool, error) {
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return false, ErrWrongBoolType
}
