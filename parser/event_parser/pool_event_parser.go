package event_parser

import (
	"github.com/ethereum/go-ethereum/common"
)

type PoolEventParser struct {
	Topic               common.Hash
	PossibleProtocolIds []int
	ethLogUnpacker      EthLogUnpacker
}
