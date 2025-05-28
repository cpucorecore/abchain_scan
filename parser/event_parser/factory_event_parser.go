package event_parser

import (
	"github.com/ethereum/go-ethereum/common"
)

type FactoryEventParser struct {
	Topic                    common.Hash
	PossibleFactoryAddresses map[common.Address]struct{}
	LogUnpacker              EthLogUnpacker
}
