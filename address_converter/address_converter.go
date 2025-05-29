package address_converter

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
)

const (
	newtonPrefix = "NEW"
)

var (
	chainIdBytes = []byte{0x03, 0xf4} // chainId 1012 hex
)

func EthAddr2ABChainAddrStr(ethAddress common.Address) string {
	data := append(chainIdBytes, ethAddress.Bytes()...)
	encoded := base58.CheckEncode(data, 0)
	newAddress := newtonPrefix + encoded
	return newAddress
}

func EthAddrStr2ABChainAddrStr(ethAddressHex string) string {
	address := common.HexToAddress(ethAddressHex)
	data := append(chainIdBytes, address.Bytes()...)
	encoded := base58.CheckEncode(data, 0)
	newAddress := newtonPrefix + encoded
	return newAddress
}
