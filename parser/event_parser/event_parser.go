package event_parser

import (
	"abchain_scan/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type EventParser interface {
	Parse(ethLog *ethtypes.Log) (types.Event, error)
}
