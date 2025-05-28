package erc20

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"log"
	"strings"
)

const (
	AbiJson = `[{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
)

var (
	Abi *abi.ABI
)

func init() {
	abiObj, err := abi.JSON(strings.NewReader(AbiJson))
	if err != nil {
		log.Fatalf("Failed to parse BEP20 ABI: %v", err)
	}
	Abi = &abiObj
}
