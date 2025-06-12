package parser

import (
	"abchain_scan/log"
	"abchain_scan/parser/event_parser"
	"abchain_scan/types"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrParserNotFound = errors.New("parser not found")
)

type TopicRouter interface {
	Parse(ethLog *ethtypes.Log) (types.Event, error)
}

type topicRouter struct {
	topic2EventParser map[common.Hash]event_parser.EventParser
}

func NewTopicRouter() TopicRouter {
	return &topicRouter{
		topic2EventParser: event_parser.Topic2EventParser,
	}
}

func (p *topicRouter) Parse(ethLog *ethtypes.Log) (types.Event, error) {
	log.Logger.Info(fmt.Sprintf("%s", ethLog.Topics[0].String()))
	eventParser, ok := p.topic2EventParser[ethLog.Topics[0]]
	if !ok {
		return nil, ErrParserNotFound
	}

	return eventParser.Parse(ethLog)
}
