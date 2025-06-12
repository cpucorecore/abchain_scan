package abi

import (
	"abchain_scan/abi/aerodrome"
	pancakev2 "abchain_scan/abi/pancake/v2"
	pancakev3 "abchain_scan/abi/pancake/v3"
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
	mapTopicToProtocolId(uniswapv2.PairCreatedTopic0, types.ProtocolIdUniswapV2)
	mapTopicToProtocolId(uniswapv2.SwapTopic0, types.ProtocolIdUniswapV2)
	mapTopicToProtocolId(uniswapv2.SyncTopic0, types.ProtocolIdUniswapV2)
	mapTopicToProtocolId(uniswapv2.BurnTopic0, types.ProtocolIdUniswapV2)
	mapTopicToProtocolId(uniswapv2.MintTopic0, types.ProtocolIdUniswapV2)

	mapTopicToProtocolId(uniswapv3.PoolCreatedTopic0, types.ProtocolIdUniswapV3)
	mapTopicToProtocolId(uniswapv3.SwapTopic0, types.ProtocolIdUniswapV3)
	mapTopicToProtocolId(uniswapv3.MintTopic0, types.ProtocolIdUniswapV3)
	mapTopicToProtocolId(uniswapv3.BurnTopic0, types.ProtocolIdUniswapV3)

	mapTopicToProtocolId(pancakev2.PairCreatedTopic0, types.ProtocolIdPancakeV2)
	mapTopicToProtocolId(pancakev2.SwapTopic0, types.ProtocolIdPancakeV2)
	mapTopicToProtocolId(pancakev2.SyncTopic0, types.ProtocolIdPancakeV2)
	mapTopicToProtocolId(pancakev2.BurnTopic0, types.ProtocolIdPancakeV2)
	mapTopicToProtocolId(pancakev2.MintTopic0, types.ProtocolIdPancakeV2)

	mapTopicToProtocolId(pancakev3.PoolCreatedTopic0, types.ProtocolIdPancakeV3)
	mapTopicToProtocolId(pancakev3.SwapTopic0, types.ProtocolIdPancakeV3)
	mapTopicToProtocolId(pancakev3.MintTopic0, types.ProtocolIdPancakeV3)
	mapTopicToProtocolId(pancakev3.BurnTopic0, types.ProtocolIdPancakeV3)

	mapTopicToProtocolId(aerodrome.PoolCreatedTopic0, types.ProtocolIdAerodrome)
	mapTopicToProtocolId(aerodrome.SwapTopic0, types.ProtocolIdAerodrome)
	mapTopicToProtocolId(aerodrome.SyncTopic0, types.ProtocolIdAerodrome)
	mapTopicToProtocolId(aerodrome.BurnTopic0, types.ProtocolIdAerodrome)
	mapTopicToProtocolId(aerodrome.MintTopic0, types.ProtocolIdAerodrome)

	FactoryAddress2ProtocolId[uniswapv2.FactoryAddress] = types.ProtocolIdUniswapV2
	FactoryAddress2ProtocolId[uniswapv3.FactoryAddress] = types.ProtocolIdUniswapV3
	FactoryAddress2ProtocolId[pancakev2.FactoryAddress] = types.ProtocolIdPancakeV2
	FactoryAddress2ProtocolId[pancakev3.FactoryAddress] = types.ProtocolIdPancakeV3
	FactoryAddress2ProtocolId[aerodrome.FactoryAddress] = types.ProtocolIdAerodrome

	mapTopicToFactoryAddress(uniswapv2.PairCreatedTopic0, uniswapv2.FactoryAddress)
	mapTopicToFactoryAddress(uniswapv3.PoolCreatedTopic0, uniswapv3.FactoryAddress)
	mapTopicToFactoryAddress(pancakev2.PairCreatedTopic0, pancakev2.FactoryAddress)
	mapTopicToFactoryAddress(pancakev3.PoolCreatedTopic0, pancakev3.FactoryAddress)
	mapTopicToFactoryAddress(aerodrome.PoolCreatedTopic0, aerodrome.FactoryAddress)
}
