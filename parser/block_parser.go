package parser

import (
	"abchain_scan/cache"
	"abchain_scan/config"
	"abchain_scan/log"
	"abchain_scan/metrics"
	"abchain_scan/sequencer"
	"abchain_scan/service"
	"abchain_scan/types"
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
	"sync"
	"time"
)

type BlockParser interface {
	Start(*sync.WaitGroup)
	Stop()
	ParseBlockAsync(bw *types.ParseBlockContext)
}

type blockParser struct {
	inputQueue   chan *types.ParseBlockContext
	workPool     *ants.Pool
	cache        cache.BlockCache
	sequencer    sequencer.Sequencer
	outputQueue  chan *types.ParseBlockContext
	priceService service.PriceService
	pairService  service.PairService
	topicRouter  TopicRouter
	kafkaSender  service.KafkaSender
	dbService    service.DBService
	parseTxPool  *ants.Pool
}

func NewBlockParser(
	cache cache.BlockCache,
	sequencer sequencer.Sequencer,
	priceService service.PriceService,
	pairService service.PairService,
	topicRouter TopicRouter,
	kafkaSender service.KafkaSender,
	dbService service.DBService,
) BlockParser {
	workPool, err := ants.NewPool(config.G.BlockHandler.PoolSize)
	if err != nil {
		log.Logger.Fatal("ants.NewPool err", zap.Error(err))
	}

	parseTxPool, err := ants.NewPool(config.G.BlockHandler.ParseTxPoolSize)
	if err != nil {
		log.Logger.Fatal("ants.NewPool err", zap.Error(err))
	}

	return &blockParser{
		inputQueue:   make(chan *types.ParseBlockContext, config.G.BlockHandler.QueueSize),
		workPool:     workPool,
		cache:        cache,
		sequencer:    sequencer,
		outputQueue:  make(chan *types.ParseBlockContext, config.G.BlockHandler.QueueSize),
		priceService: priceService,
		pairService:  pairService,
		topicRouter:  topicRouter,
		kafkaSender:  kafkaSender,
		dbService:    dbService,
		parseTxPool:  parseTxPool,
	}
}

func (p *blockParser) Commit(x sequencer.Sequenceable) {
	p.outputQueue <- x.(*types.ParseBlockContext)
}

func (p *blockParser) Start(waitGroup *sync.WaitGroup) {
	p.startHandleBlockResult(waitGroup)

	go func() {
		wg := &sync.WaitGroup{}
	tagFor:
		for {
			select {
			case pbc, ok := <-p.inputQueue:
				if !ok {
					log.Logger.Info("block handler inputQueue is closed")
					break tagFor
				}

				wg.Add(1)
				p.workPool.Submit(func() {
					defer wg.Done()
					p.parseBlock(pbc)
				})
			}
		}

		wg.Wait()
		log.Logger.Info("all block parse task finish")
		p.doStop()
	}()
}

func (p *blockParser) Stop() {
	close(p.inputQueue)
}

func (p *blockParser) ParseBlockAsync(bw *types.ParseBlockContext) {
	p.inputQueue <- bw
}

func (p *blockParser) waitForNativeTokenPrice(blockNumber *big.Int) decimal.Decimal {
	for {
		bnbPrice, err := p.priceService.GetNativeTokenPrice(blockNumber)
		if err != nil {
			log.Logger.Error("get price err", zap.Error(err), zap.Any("blockNumber", blockNumber))
			time.Sleep(time.Millisecond * 100)
			continue
		}
		return bnbPrice
	}
}

func collectNewPairAndTokens(br *types.BlockResult, pairWrap *types.PairWrap) {
	if pairWrap.NewPair {
		br.NewPairs[pairWrap.Pair.Address] = pairWrap.Pair
	}

	if pairWrap.NewToken0 {
		token0 := pairWrap.Pair.Token0
		br.NewTokens[token0.Address] = token0
	}

	if pairWrap.NewToken1 {
		token1 := pairWrap.Pair.Token1
		br.NewTokens[token1.Address] = token1
	}
}

type TxResultAndPairWrap struct {
	TxResult  *types.TxResult
	PairWraps []*types.PairWrap
}

func (p *blockParser) parseTxReceipt(pbc *types.ParseBlockContext, txReceipt *ethtypes.Receipt) *TxResultAndPairWrap {
	txSender, err := pbc.GetTxSender(txReceipt.TransactionIndex)
	if err != nil {
		log.Logger.Fatal("Err: get tx sender err",
			zap.Error(err),
			zap.Any("pbc", pbc),
			zap.Any("txReceipt", txReceipt),
		)
	}

	tr := types.NewTxResult(txSender)
	pairWraps := make([]*types.PairWrap, 0, len(txReceipt.Logs))
	for _, ethLog := range txReceipt.Logs {
		if len(ethLog.Topics) == 0 {
			continue
		}

		event, parseErr := p.topicRouter.Parse(ethLog)
		if parseErr != nil {
			continue
		}

		pairWrap := p.getPairByEvent(event)
		if pairWrap.Pair.Filtered {
			continue
		}

		pairWraps = append(pairWraps, pairWrap)
		event.SetPair(pairWrap.Pair)
		event.SetBlockTime(pbc.HeightTime.Time)
		tr.AddEvent(event)
	}

	return &TxResultAndPairWrap{
		TxResult:  tr,
		PairWraps: pairWraps,
	}
}

func (p *blockParser) parseBlock(pbc *types.ParseBlockContext) {
	pbc.NativeTokenPrice = p.waitForNativeTokenPrice(pbc.HeightTime.HeightBigInt)

	now := time.Now()
	br := types.NewBlockResult(pbc.HeightTime.Height, pbc.HeightTime.Timestamp, pbc.NativeTokenPrice)

	wg := &sync.WaitGroup{}
	results := make([]*TxResultAndPairWrap, len(pbc.BlockReceipts))
	for i, txReceipt := range pbc.BlockReceipts {
		if txReceipt.Status != 1 {
			continue
		}

		wg.Add(1)
		p.parseTxPool.Submit(func() {
			defer wg.Done()
			result := p.parseTxReceipt(pbc, txReceipt)
			results[i] = result
		})
	}
	wg.Wait()

	for _, result := range results {
		if result == nil {
			continue
		}

		for _, pairWrap := range result.PairWraps {
			collectNewPairAndTokens(br, pairWrap)
		}
		br.AddTxResult(result.TxResult)
	}

	duration := time.Since(now)
	metrics.ParseBlockDurationMs.Observe(float64(duration.Milliseconds()))
	log.Logger.Info(fmt.Sprintf("parse block %d duration %dms", pbc.HeightTime.HeightBigInt, duration.Milliseconds()))

	pbc.BlockResult = br
	p.sequencer.CommitWithSequence(pbc, p)
}

func (p *blockParser) getPairByEvent(event types.Event) *types.PairWrap {
	if event.CanGetPair() {
		pair := event.GetPair()
		if pair.Filtered {
			p.pairService.SetPair(pair)
			return &types.PairWrap{
				Pair:      pair,
				NewPair:   false,
				NewToken0: false,
				NewToken1: false,
			}
		}

		return p.pairService.GetPairTokens(pair)
	}

	return p.pairService.GetPair(event.GetPairAddress(), event.GetPossibleProtocolIds())
}

func (p *blockParser) commitBlockResult(blockResult *types.BlockResult) {
	blockInfo := blockResult.GetKafkaMessage()

	now := time.Now()
	err := p.dbService.AddTokens(blockInfo.NewTokens)
	if err != nil {
		log.Logger.Fatal("add tokens err", zap.Any("height", blockInfo.Height), zap.Error(err))
	}

	err = p.dbService.AddPairs(blockInfo.NewPairs)
	if err != nil {
		log.Logger.Fatal("add pairs err", zap.Any("height", blockInfo.Height), zap.Error(err))
	}

	err = p.dbService.AddTxs(blockInfo.Txs)
	if err != nil {
		log.Logger.Fatal("add txs err", zap.Any("height", blockInfo.Height), zap.Error(err))
	}

	duration := time.Since(now)
	metrics.DbOperationDurationMs.Observe(float64(duration.Milliseconds()))
	log.Logger.Info("db operation duration",
		zap.Uint64("block", blockResult.Height),
		zap.Float64("duration", duration.Seconds()),
		zap.String("price", blockInfo.NativeTokenPrice),
		zap.Int("new tokens", len(blockInfo.NewTokens)),
		zap.Int("new pairs", len(blockInfo.NewPairs)),
		zap.Int("txs", len(blockInfo.Txs)))

	err = p.kafkaSender.Send(blockInfo)
	if err != nil {
		log.Logger.Fatal("kafka send msg err", zap.Error(err), zap.Any("block", blockResult.Height))
	}

	p.cache.SetFinishedBlock(blockResult.Height)
	metrics.CurrentHeight.Set(float64(blockResult.Height))
	metrics.TxCntByBlock.Set(float64(len(blockInfo.Txs)))
}

func (p *blockParser) startHandleBlockResult(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {
			blockContext, ok := <-p.outputQueue
			if !ok {
				log.Logger.Info("commitBlockResultOld - output queue closed")
				return
			}

			p.commitBlockResult(blockContext.BlockResult)
		}
	}()
}

func (p *blockParser) doStop() {
	p.workPool.Release()
	close(p.outputQueue)
}
