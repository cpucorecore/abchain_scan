package repository

import (
	"abchain_scan/repository/orm"
	"abchain_scan/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

func prepareTxTest() *TxRepository {
	dsn := "host=localhost user=postgres password=12345678 dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	return NewTxRepository(db)
}

func cleanupTxTest(txRepository *TxRepository, Ids ...string) {
	for _, id := range Ids {
		txRepository.DeleteById(id)
	}
}

func TestTxRepository_Create(t *testing.T) {
	txRepository := prepareTxTest()
	tx := &orm.Tx{
		TxHash:        "0xa1",
		Event:         "buy",
		Token0Amount:  decimal.NewFromFloat(0.001),
		Token1Amount:  decimal.NewFromFloat(0.002),
		Maker:         "0xa1",
		Token0Address: "0xa1",
		Token1Address: "0xa1",
		AmountUsd:     decimal.NewFromFloat(0.003),
		PriceUsd:      decimal.NewFromFloat(0.004),
		Block:         1,
		BlockAt:       time.Now(),
		BlockIndex:    1,
		TxIndex:       1,
		PairAddress:   "0xa1",
		Program:       types.ProtocolNameUniswapV2,
	}
	createErr := txRepository.Create(tx)
	require.NoError(t, createErr)

	pairQueried, getErr := txRepository.GetByUniqIndex(tx.Token0Address, tx.Block, tx.BlockIndex, tx.TxIndex)
	require.Nil(t, getErr)
	require.True(t, tx.Equal(pairQueried))
	cleanupTxTest(txRepository, pairQueried.Id.String())
}

func TestTxRepository_CreateDup(t *testing.T) {
	txRepository := prepareTxTest()
	txes := []*orm.Tx{
		{
			TxHash:        "0xa1",
			Event:         "buy",
			Token0Amount:  decimal.NewFromFloat(0.001),
			Token1Amount:  decimal.NewFromFloat(0.002),
			Maker:         "0xa1",
			Token0Address: "0xa1",
			Token1Address: "0xa1",
			AmountUsd:     decimal.NewFromFloat(0.003),
			PriceUsd:      decimal.NewFromFloat(0.004),
			Block:         1,
			BlockAt:       time.Now(),
			BlockIndex:    1,
			TxIndex:       1,
			PairAddress:   "0xa1",
			Program:       types.ProtocolNameUniswapV2,
		},
		{
			TxHash:        "0xa2",
			Event:         "buy",
			Token0Amount:  decimal.NewFromFloat(0.001),
			Token1Amount:  decimal.NewFromFloat(0.002),
			Maker:         "0xa2",
			Token0Address: "0xa2",
			Token1Address: "0xa2",
			AmountUsd:     decimal.NewFromFloat(0.003),
			PriceUsd:      decimal.NewFromFloat(0.004),
			Block:         2,
			BlockAt:       time.Now(),
			BlockIndex:    2,
			TxIndex:       2,
			PairAddress:   "0xa2",
			Program:       types.ProtocolNameUniswapV2,
		},
		{
			TxHash:        "0xa3",
			Event:         "buy",
			Token0Amount:  decimal.NewFromFloat(0.001),
			Token1Amount:  decimal.NewFromFloat(0.002),
			Maker:         "0xa3",
			Token0Address: "0xa3",
			Token1Address: "0xa3",
			AmountUsd:     decimal.NewFromFloat(0.003),
			PriceUsd:      decimal.NewFromFloat(0.004),
			Block:         3,
			BlockAt:       time.Now(),
			BlockIndex:    3,
			TxIndex:       3,
			PairAddress:   "0xa3",
			Program:       types.ProtocolNameUniswapV2,
		},
	}

	createErr := txRepository.Create(txes[0])
	require.NoError(t, createErr)

	createBatchErr := txRepository.CreateBatch(txes, "token0_address", "block", "block_index", "tx_index")
	require.NoError(t, createBatchErr)

	txIds := make([]string, 0, 3)
	for i, tx := range txes {
		pairQueried, err := txRepository.GetByUniqIndex(tx.Token0Address, tx.Block, tx.BlockIndex, tx.TxIndex)
		require.Nil(t, err)
		require.True(t, txes[i].Equal(pairQueried))
		txIds = append(txIds, pairQueried.Id.String())
	}

	defer cleanupTxTest(txRepository, txIds...)
}
