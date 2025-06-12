package abi

import (
	uniswapv2 "abchain_scan/abi/uniswap/v2"
	uniswapv3 "abchain_scan/abi/uniswap/v3"
	"abchain_scan/types"
	"github.com/ethereum/go-ethereum/common"
)

var Topic2ProtocolIds = map[common.Hash][]int{}
var FactoryAddress2ProtocolId = map[common.Address]int{}
var Topic2FactoryAddresses = map[common.Hash]map[common.Address]struct{}{}

func mapTopicToProtocolId(topic common.Hash, protocolId int) {
	protocolIds, ok := Topic2ProtocolIds[topic]
	if !ok {
		protocolIds = []int{}
	}
	protocolIds = append(protocolIds, protocolId)
	Topic2ProtocolIds[topic] = protocolIds
}

func mapTopicToFactoryAddress(topic common.Hash, factoryAddress common.Address) {
	factoryAddresses, ok := Topic2FactoryAddresses[topic]
	if !ok {
		factoryAddresses = make(map[common.Address]struct{})
	}
	factoryAddresses[factoryAddress] = struct{}{}
	Topic2FactoryAddresses[topic] = factoryAddresses
}

func init() {
	mapTopicToProtocolId(uniswapv2.PairCreatedTopic0, types.ProtocolIdNewSwap)
	mapTopicToProtocolId(uniswapv2.SwapTopic0, types.ProtocolIdNewSwap)
	mapTopicToProtocolId(uniswapv2.SyncTopic0, types.ProtocolIdNewSwap)
	mapTopicToProtocolId(uniswapv2.BurnTopic0, types.ProtocolIdNewSwap)
	mapTopicToProtocolId(uniswapv2.MintTopic0, types.ProtocolIdNewSwap)

	mapTopicToProtocolId(uniswapv3.PoolCreatedTopic0, types.ProtocolIdUniswapV3)
	mapTopicToProtocolId(uniswapv3.SwapTopic0, types.ProtocolIdUniswapV3)
	mapTopicToProtocolId(uniswapv3.MintTopic0, types.ProtocolIdUniswapV3)
	mapTopicToProtocolId(uniswapv3.BurnTopic0, types.ProtocolIdUniswapV3)

	FactoryAddress2ProtocolId[uniswapv2.FactoryAddress] = types.ProtocolIdNewSwap
	FactoryAddress2ProtocolId[uniswapv3.FactoryAddress] = types.ProtocolIdUniswapV3

	mapTopicToFactoryAddress(uniswapv2.PairCreatedTopic0, uniswapv2.FactoryAddress)
	mapTopicToFactoryAddress(uniswapv3.PoolCreatedTopic0, uniswapv3.FactoryAddress)
}
