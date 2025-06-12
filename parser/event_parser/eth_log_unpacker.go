package event_parser

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrWrongTopicLen      = errors.New("wrong topic length")
	ErrWrongDataUnpackLen = errors.New("wrong data unpack length")
)

type EthLogUnpacker struct {
	AbiEvent      *abi.Event
	TopicLen      int
	DataUnpackLen int
}

func (p *EthLogUnpacker) Unpack(ethLog *ethtypes.Log) ([]interface{}, error) {
	if len(ethLog.Topics) != p.TopicLen {
		return nil, ErrWrongTopicLen
	}

	eventInput, err := p.AbiEvent.Inputs.Unpack(ethLog.Data)
	if err != nil {
		return nil, err
	}

	if len(eventInput) != p.DataUnpackLen {
		return nil, ErrWrongDataUnpackLen
	}

	return eventInput, nil
}
