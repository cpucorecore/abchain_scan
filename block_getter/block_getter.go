package block_getter

import (
	"abchain_scan/cache"
	"abchain_scan/config"
	"abchain_scan/http_client"
	"abchain_scan/log"
	"abchain_scan/metrics"
	"abchain_scan/sequencer"
	"abchain_scan/types"
	"context"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"math/big"
	"sync"
	"time"
)

type BlockGetter interface {
	Start()
	GetStartBlockNumber(startBlockNumber uint64) uint64
	StartDispatch(startBlockNumber uint64)
	Stop()
	GetBlockAsync(blockNumber uint64)
	Next() *types.ParseBlockContext
}

type blockGetter struct {
	ctx                  context.Context
	ethClient            *ethclient.Client
	wsEthClient          *ethclient.Client
	inputQueue           chan uint64
	outputBuffer         chan *types.ParseBlockContext
	workPool             *ants.Pool
	cache                cache.BlockCache
	stopped              SafeVar[bool]
	blockHeaderChan      chan *ethtypes.Header
	blockSequencer       sequencer.Sequencer
	headerHeight         SafeVar[uint64]
	retryParams          *config.RetryParams
	getTxReceiptWorkPool *ants.Pool
}

func NewBlockGetter(ethClient *ethclient.Client,
	wsEthClient *ethclient.Client,
	cache cache.BlockCache,
	blockSequencer sequencer.Sequencer,
	retryParams *config.RetryParams,
) BlockGetter {
	workPool, err := ants.NewPool(config.G.BlockGetter.PoolSize)
	if err != nil {
		log.Logger.Fatal("ants pool(BlockGetter) init err", zap.Error(err))
	}

	getTxReceiptWorkPool, err := ants.NewPool(config.G.BlockGetter.GetTxReceiptPoolSize)
	if err != nil {
		log.Logger.Fatal("ants pool(GetTxReceipt) init err", zap.Error(err))
	}

	return &blockGetter{
		ctx:                  context.Background(),
		ethClient:            ethClient,
		wsEthClient:          wsEthClient,
		inputQueue:           make(chan uint64, config.G.BlockGetter.QueueSize),
		outputBuffer:         make(chan *types.ParseBlockContext, 10),
		workPool:             workPool,
		cache:                cache,
		blockHeaderChan:      make(chan *ethtypes.Header, 100),
		blockSequencer:       blockSequencer,
		retryParams:          retryParams,
		getTxReceiptWorkPool: getTxReceiptWorkPool,
	}
}

func (bg *blockGetter) Commit(x sequencer.Sequenceable) {
	bg.outputBuffer <- x.(*types.ParseBlockContext)
}

func (bg *blockGetter) getTxReceiptRetry(txHash common.Hash) (*ethtypes.Receipt, error) {
	return retry.DoWithData(func() (*ethtypes.Receipt, error) {
		txReceipt, err := bg.ethClient.TransactionReceipt(bg.ctx, txHash)
		if err != nil {
			log.Logger.Error("TransactionReceipt() err", zap.String("txHash", txHash.String()), zap.Error(err))
			return nil, err
		}
		return txReceipt, err
	}, bg.retryParams.Attempts, bg.retryParams.Delay)
}

func (bg *blockGetter) getBlock(blockNumber uint64) (*types.ParseBlockContext, error) {
	var (
		block       *ethtypes.Block
		getBlockErr error
	)

	now := time.Now()
	block, getBlockErr = bg.ethClient.BlockByNumber(bg.ctx, big.NewInt(int64(blockNumber)))
	if getBlockErr == nil {
		duration := time.Since(now)
		metrics.GetBlockDurationMs.Observe(float64(duration.Milliseconds()))
	} else {
		log.Logger.Error("BlockByNumber() err", zap.Uint64("blockNumber", blockNumber), zap.Error(getBlockErr))
		return nil, getBlockErr
	}

	now = time.Now()
	wg := &sync.WaitGroup{}
	var getTxReceiptErr error
	mu := &sync.Mutex{}
	blockReceipts := make([]*ethtypes.Receipt, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		wg.Add(1)
		bg.getTxReceiptWorkPool.Submit(func() {
			defer wg.Done()

			txReceipt, err := bg.getTxReceiptRetry(tx.Hash())
			if err != nil {
				log.Logger.Error("TransactionReceipt() err", zap.Uint64("blockNumber", blockNumber), zap.Any("tx_hash", tx.Hash()), zap.Error(err))
				mu.Lock()
				getTxReceiptErr = err
				mu.Unlock()
				return
			}
			blockReceipts[i] = txReceipt
		})
	}
	wg.Wait()

	if getTxReceiptErr != nil {
		log.Logger.Error("Get BlockReceipts err", zap.Uint64("blockNumber", blockNumber), zap.Error(getTxReceiptErr))
		return nil, getTxReceiptErr
	}

	duration := time.Since(now)
	if duration.Seconds() > 1 {
		log.Logger.Info("Get BlockReceipts", zap.Uint64("blockNumber", blockNumber), zap.Duration("duration", duration))
	}
	metrics.GetBlockReceiptsDurationMs.Observe(float64(duration.Milliseconds()))

	if getBlockErr != nil {
		return nil, getBlockErr
	}

	metrics.BlockDelay.Observe(time.Now().Sub(time.Unix((int64)(block.Time()), 0)).Seconds())

	transactions := block.Transactions()
	return &types.ParseBlockContext{
		Block:           block,
		Transactions:    transactions,
		TransactionsLen: uint(len(transactions)),
		BlockReceipts:   blockReceipts,
		HeightTime:      types.GetBlockHeightTime(block.Header()),
		TxSenders:       make([]*common.Address, block.Transactions().Len()),
	}, nil
}

func (bg *blockGetter) getBlockWithRetry(blockNumber uint64) (*types.ParseBlockContext, error) {
	return retry.DoWithData(func() (*types.ParseBlockContext, error) {
		return bg.getBlock(blockNumber)
	}, bg.retryParams.Attempts, bg.retryParams.Delay)
}

func (bg *blockGetter) GetBlockAsync(blockNumber uint64) {
	bg.inputQueue <- blockNumber
}

func (bg *blockGetter) Next() *types.ParseBlockContext {
	return <-bg.outputBuffer
}

func (bg *blockGetter) Start() {
	go func() {
		wg := &sync.WaitGroup{}
	tagFor:
		for {
			select {
			case blockNumber, ok := <-bg.inputQueue:
				if !ok {
					log.Logger.Info("block inputQueue is closed")
					break tagFor
				}

				wg.Add(1)
				bg.workPool.Submit(func() {
					defer wg.Done()

					log.Logger.Info("get block start", zap.Uint64("block_number", blockNumber))
					bw, err := bg.getBlockWithRetry(blockNumber)
					if err != nil {
						log.Logger.Error("get block err", zap.Uint64("blockNumber", blockNumber), zap.Error(err))
						return
					}

					log.Logger.Info("get block success", zap.Uint64("blockNumber", blockNumber))
					metrics.BlockQueueSize.Set(float64(len(bg.outputBuffer)))
					bg.blockSequencer.CommitWithSequence(bw, bg)
				})
			}
		}

		wg.Wait()
		log.Logger.Info("all block getter task finish")
		close(bg.outputBuffer)
	}()
}

func (bg *blockGetter) GetStartBlockNumber(startBlockNumber uint64) uint64 {
	if startBlockNumber != 0 {
		return startBlockNumber
	}

	finishedBlock := bg.cache.GetFinishedBlock()
	if finishedBlock != 0 {
		return finishedBlock + 1
	}

	newestBlockNumber, err := http_client.GetLatestBlockNumber(config.G.Chain.Endpoint)
	if err != nil {
		log.Logger.Fatal("ethClient.BlockNumber() err", zap.Error(err))
	}

	return newestBlockNumber
}

func (bg *blockGetter) setHeaderHeight(headerHeight uint64) {
	if headerHeight > bg.headerHeight.Get() {
		bg.headerHeight.Set(headerHeight)
	}
}

func (bg *blockGetter) getHeaderHeight() uint64 {
	return bg.headerHeight.Get()
}

func (bg *blockGetter) subscribeNewHead() (ethereum.Subscription, <-chan error, error) {
	sub, err := bg.wsEthClient.SubscribeNewHead(bg.ctx, bg.blockHeaderChan)
	if err != nil {
		return nil, nil, err
	}
	return sub, sub.Err(), nil
}

func (bg *blockGetter) reconnectWithBackoff() (ethereum.Subscription, <-chan error) {
	retryDelay := time.Second * 1
	maxRetryDelay := time.Second * 10

	for {
		sub, errChan, err := bg.subscribeNewHead()
		if err == nil {
			log.Logger.Info("WebSocket reconnected successfully")
			return sub, errChan
		}

		log.Logger.Error("WebSocket reconnect failed",
			zap.Error(err),
			zap.Duration("nextRetry", retryDelay),
		)
		time.Sleep(retryDelay)

		retryDelay *= 2
		if retryDelay > maxRetryDelay {
			retryDelay = maxRetryDelay
		}
	}
}

func (bg *blockGetter) startSubscribeNewHead() {
	headerHeight, err := http_client.GetLatestBlockNumber(config.G.Chain.Endpoint)
	if err != nil {
		log.Logger.Fatal("HeightBigInt() err", zap.Error(err))
	}
	bg.setHeaderHeight(headerHeight)

	sub, errChan, err := bg.subscribeNewHead()
	if err != nil {
		log.Logger.Fatal("subscribeNewHead() err", zap.Error(err))
	}

	go func() {
		noBlockTimeout := time.NewTimer(10 * time.Second)
		defer noBlockTimeout.Stop()

		resetConnection := func() {
			noBlockTimeout.Stop()
			select {
			case <-noBlockTimeout.C:
			default:
			}
			sub.Unsubscribe()

			sub, errChan = bg.reconnectWithBackoff()
			noBlockTimeout.Reset(10 * time.Second)
		}

		for {
			select {
			case err = <-errChan:
				log.Logger.Error("WebSocket error", zap.Error(err))
				resetConnection()
			case blockHeader := <-bg.blockHeaderChan:
				height := blockHeader.Number.Uint64()
				log.Logger.Info("New block", zap.Uint64("height", height))
				bg.setHeaderHeight(height)
				metrics.NewestHeight.Set(float64(height))

				noBlockTimeout.Stop()
				select {
				case <-noBlockTimeout.C:
				default:
				}
				noBlockTimeout.Reset(10 * time.Second)
			case <-noBlockTimeout.C:
				log.Logger.Warn("No new blocks for 10s, reconnect WebSocket")
				resetConnection()
			}
		}
	}()
}

func (bg *blockGetter) dispatchRange(from, to uint64) (stopped bool, nextBlock uint64) {
	for i := from; i <= to; i++ {
		if bg.isStopped() {
			return true, i
		}
		bg.GetBlockAsync(i)
	}
	return false, 0
}

func (bg *blockGetter) StartDispatch(startBlockNumber uint64) {
	bg.startSubscribeNewHead()

	go func() {
		cur := startBlockNumber
		for {
			headerHeight := bg.getHeaderHeight()
			if headerHeight < cur {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			stopped, nextBlockHeight := bg.dispatchRange(cur, headerHeight)
			if stopped {
				log.Logger.Info("dispatch interrupted", zap.Uint64("nextBlockHeight", nextBlockHeight))
				bg.doStop()
				return
			}

			cur = headerHeight + 1
		}
	}()
}

func (bg *blockGetter) Stop() {
	bg.stopped.Set(true)
}

func (bg *blockGetter) isStopped() bool {
	return bg.stopped.Get()
}

func (bg *blockGetter) doStop() {
	close(bg.inputQueue)
}
