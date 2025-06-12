package types

import (
	"abchain_scan/log"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type TxPairEvent struct {
	UniswapV2 []Event
	UniswapV3 []Event
	PancakeV2 []Event
	PancakeV3 []Event
	Aerodrome []Event
}

func (tpe *TxPairEvent) AddEvent(event Event) {
	switch event.GetProtocolId() {
	case ProtocolIdUniswapV2:
		if tpe.UniswapV2 == nil {
			tpe.UniswapV2 = make([]Event, 0, 10)
		}
		tpe.UniswapV2 = append(tpe.UniswapV2, event)
	case ProtocolIdUniswapV3:
		if tpe.UniswapV3 == nil {
			tpe.UniswapV3 = make([]Event, 0, 10)
		}
		tpe.UniswapV3 = append(tpe.UniswapV3, event)
	case ProtocolIdPancakeV2:
		if tpe.PancakeV2 == nil {
			tpe.PancakeV2 = make([]Event, 0, 10)
		}
		tpe.PancakeV2 = append(tpe.PancakeV2, event)
	case ProtocolIdPancakeV3:
		if tpe.PancakeV3 == nil {
			tpe.PancakeV3 = make([]Event, 0, 10)
		}
		tpe.PancakeV3 = append(tpe.PancakeV3, event)
	case ProtocolIdAerodrome:
		if tpe.Aerodrome == nil {
			tpe.Aerodrome = make([]Event, 0, 10)
		}
		tpe.Aerodrome = append(tpe.Aerodrome, event)
	}
}

func (tpe *TxPairEvent) LinkEvents() {
	tpe.linkEventByProtocol(tpe.UniswapV2)
	tpe.linkEventByProtocol(tpe.UniswapV3)
	tpe.linkEventByProtocol(tpe.PancakeV2)
	tpe.linkEventByProtocol(tpe.PancakeV3)
	tpe.linkEventByProtocol(tpe.Aerodrome)
}

func LinkPairCreatedEventAndMintEvent(pairCreatedEvents, mintEvents []Event) {
	mintEventsLen := len(mintEvents)
	for i, pairCreatedEvent := range pairCreatedEvents {
		if i < mintEventsLen {
			pairCreatedEvent.LinkEvent(mintEvents[i])
		} else {
			log.Logger.Info("Waring: pair have no related mint event", zap.Any("pairCreatedEvent", pairCreatedEvent))
		}
	}
}

func (tpe *TxPairEvent) linkEventByProtocol(events []Event) {
	mintEvents := make([]Event, 0, 10)
	pairCreatedEvents := make([]Event, 0, 10)
	for _, event := range events {
		if event.IsMint() {
			mintEvents = append(mintEvents, event)
		} else if event.IsCreatePair() {
			pairCreatedEvents = append(pairCreatedEvents, event)
		}
	}
	LinkPairCreatedEventAndMintEvent(pairCreatedEvents, mintEvents)
}

type TxResult struct {
	Maker                   common.Address
	PairCreatedEvents       []Event
	PairAddress2TxPairEvent map[common.Address]*TxPairEvent
}

func NewTxResult(maker common.Address) *TxResult {
	return &TxResult{
		Maker:                   maker,
		PairCreatedEvents:       make([]Event, 0, 10),
		PairAddress2TxPairEvent: make(map[common.Address]*TxPairEvent),
	}
}

func (tr *TxResult) AddEvent(event Event) {
	event.SetMaker(tr.Maker)
	if event.IsCreatePair() {
		tr.PairCreatedEvents = append(tr.PairCreatedEvents, event)
	}

	pairAddress := event.GetPairAddress()
	txPairEvent, ok := tr.PairAddress2TxPairEvent[pairAddress]
	if ok {
		txPairEvent.AddEvent(event)
		return
	}

	txPairEvent = &TxPairEvent{}
	txPairEvent.AddEvent(event)
	tr.PairAddress2TxPairEvent[pairAddress] = txPairEvent
}

func (tr *TxResult) LinkEvents() {
	for _, pairEvent := range tr.PairAddress2TxPairEvent {
		pairEvent.LinkEvents()
	}
}
