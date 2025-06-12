package types

import (
	chainparams "abchain_scan/chain"
	"abchain_scan/log"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
	"time"
)

var (
	txIndexOutOfRange = errors.New("txIndex out of range")
)

type BlockHeightTime struct {
	Height       uint64
	Timestamp    uint64
	HeightBigInt *big.Int
	Time         time.Time
}

func GetBlockHeightTime(header *ethtypes.Header) *BlockHeightTime {
	return &BlockHeightTime{
		HeightBigInt: header.Number,
		Height:       header.Number.Uint64(),
		Timestamp:    header.Time,
		Time:         time.Unix((int64)(header.Time), 0).UTC(),
	}
}

type ParseBlockContext struct {
	// input
	Block            *ethtypes.Block
	Transactions     []*ethtypes.Transaction
	TransactionsLen  uint
	BlockReceipts    []*ethtypes.Receipt
	HeightTime       *BlockHeightTime
	NativeTokenPrice decimal.Decimal
	TxSenders        []*common.Address
	// output
	BlockResult *BlockResult
}

func (c *ParseBlockContext) GetSequence() uint64 {
	return c.HeightTime.Height
}

func (c *ParseBlockContext) GetTxSender(txIndex uint) (common.Address, error) {
	if c.TxSenders[txIndex] != nil {
		return *c.TxSenders[txIndex], nil
	}

	if txIndex >= c.TransactionsLen {
		log.Logger.Info("Waring: txIndex out of range",
			zap.Uint64("height", c.HeightTime.Height),
			zap.Any("transactions length", c.TransactionsLen),
			zap.Uint("txIndex", txIndex),
		)
		return ZeroAddress, txIndexOutOfRange
	}

	signer := ethtypes.MakeSigner(chainparams.ChainConfig, c.HeightTime.HeightBigInt, c.HeightTime.Timestamp)
	sender, err := ethtypes.Sender(signer, c.Transactions[txIndex])
	if err != nil {
		return ZeroAddress, err
	}

	c.TxSenders[txIndex] = &sender
	return sender, nil
}
