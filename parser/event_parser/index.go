package event_parser

import (
	"abchain_scan/abi"
	uniswapv2 "abchain_scan/abi/uniswap/v2"
	"github.com/ethereum/go-ethereum/common"
)

var (
	pairCreatedEventParser = &PairCreatedEventParser{
		FactoryEventParser: FactoryEventParser{
			Topic:                    uniswapv2.PairCreatedTopic0,
			PossibleFactoryAddresses: abi.Topic2FactoryAddresses[uniswapv2.PairCreatedTopic0],
			LogUnpacker: EthLogUnpacker{
				AbiEvent:      uniswapv2.PairCreatedEvent,
				TopicLen:      3,
				DataUnpackLen: 2,
			},
		},
	}

	burnEventParser = &BurnEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               uniswapv2.BurnTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv2.BurnTopic0],
			ethLogUnpacker: EthLogUnpacker{
				AbiEvent:      uniswapv2.BurnEvent,
				TopicLen:      3,
				DataUnpackLen: 2,
			},
		},
	}

	swapEventParser = &SwapEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               uniswapv2.SwapTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv2.SwapTopic0],
			ethLogUnpacker: EthLogUnpacker{
				AbiEvent:      uniswapv2.SwapEvent,
				TopicLen:      3,
				DataUnpackLen: 4,
			},
		},
	}

	syncEventParser = &SyncEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               uniswapv2.SyncTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv2.SyncTopic0],
			ethLogUnpacker: EthLogUnpacker{
				AbiEvent:      uniswapv2.SyncEvent,
				TopicLen:      1,
				DataUnpackLen: 2,
			},
		},
	}

	mintEventParser = &MintEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               uniswapv2.MintTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv2.MintTopic0],
			ethLogUnpacker: EthLogUnpacker{
				AbiEvent:      uniswapv2.MintEvent,
				TopicLen:      2,
				DataUnpackLen: 2,
			},
		},
	}

	Topic2EventParser = map[common.Hash]EventParser{
		uniswapv2.PairCreatedTopic0: pairCreatedEventParser,
		uniswapv2.MintTopic0:        mintEventParser,
		uniswapv2.BurnTopic0:        burnEventParser,
		uniswapv2.SwapTopic0:        swapEventParser,
		uniswapv2.SyncTopic0:        syncEventParser,
	}
)
