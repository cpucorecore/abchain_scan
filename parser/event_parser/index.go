package event_parser

import (
	"abchain_scan/abi"
	"abchain_scan/abi/aerodrome"
	pancakev3 "abchain_scan/abi/pancake/v3"
	uniswapv2 "abchain_scan/abi/uniswap/v2"
	uniswapv3 "abchain_scan/abi/uniswap/v3"
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

	poolCreatedEventParserAerodrome = &PairCreatedEventParser{
		FactoryEventParser: FactoryEventParser{
			Topic:                    aerodrome.PoolCreatedTopic0,
			PossibleFactoryAddresses: abi.Topic2FactoryAddresses[aerodrome.PoolCreatedTopic0],
			LogUnpacker: EthLogUnpacker{
				AbiEvent:      aerodrome.PoolCreatedEvent,
				TopicLen:      4,
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

	burnEventParserAerodrome = &BurnEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               aerodrome.BurnTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[aerodrome.BurnTopic0],
			ethLogUnpacker: EthLogUnpacker{
				AbiEvent:      aerodrome.BurnEvent,
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

	swapEventParserAerodrome = &SwapEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               aerodrome.SwapTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[aerodrome.SwapTopic0],
			ethLogUnpacker: EthLogUnpacker{
				AbiEvent:      aerodrome.SwapEvent,
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

	syncEventParserAerodrome = &SyncEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               aerodrome.SyncTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[aerodrome.SyncTopic0],
			ethLogUnpacker: EthLogUnpacker{
				AbiEvent:      aerodrome.SyncEvent,
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

		uniswapv3.PoolCreatedTopic0: &PoolCreatedEventParser{
			FactoryEventParser: FactoryEventParser{
				Topic:                    uniswapv3.PoolCreatedTopic0,
				PossibleFactoryAddresses: abi.Topic2FactoryAddresses[uniswapv3.PoolCreatedTopic0],
				LogUnpacker: EthLogUnpacker{
					AbiEvent:      uniswapv3.PoolCreatedEvent,
					TopicLen:      4,
					DataUnpackLen: 2,
				},
			},
		},
		uniswapv3.MintTopic0: &MintEventParserV3{
			PoolEventParser: PoolEventParser{
				Topic:               uniswapv3.MintTopic0,
				PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv3.MintTopic0],
				ethLogUnpacker: EthLogUnpacker{
					AbiEvent:      uniswapv3.MintEvent,
					TopicLen:      4,
					DataUnpackLen: 4,
				},
			},
		},
		uniswapv3.BurnTopic0: &BurnEventParserV3{
			PoolEventParser: PoolEventParser{
				Topic:               uniswapv3.BurnTopic0,
				PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv3.BurnTopic0],
				ethLogUnpacker: EthLogUnpacker{
					AbiEvent:      uniswapv3.BurnEvent,
					TopicLen:      4,
					DataUnpackLen: 3,
				},
			},
		},
		uniswapv3.SwapTopic0: &SwapEventParserV3{
			PoolEventParser: PoolEventParser{
				Topic:               uniswapv3.SwapTopic0,
				PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv3.SwapTopic0],
				ethLogUnpacker: EthLogUnpacker{
					AbiEvent:      uniswapv3.SwapEvent,
					TopicLen:      3,
					DataUnpackLen: 5,
				},
			},
		},

		pancakev3.SwapTopic0: &SwapEventParserV3{
			PoolEventParser: PoolEventParser{
				Topic:               pancakev3.SwapTopic0,
				PossibleProtocolIds: abi.Topic2ProtocolIds[pancakev3.SwapTopic0],
				ethLogUnpacker: EthLogUnpacker{
					AbiEvent:      pancakev3.SwapEvent,
					TopicLen:      3,
					DataUnpackLen: 7,
				},
			},
		},

		aerodrome.PoolCreatedTopic0: poolCreatedEventParserAerodrome,
		aerodrome.BurnTopic0:        burnEventParserAerodrome,
		aerodrome.SwapTopic0:        swapEventParserAerodrome,
		aerodrome.SyncTopic0:        syncEventParserAerodrome,
	}
)
