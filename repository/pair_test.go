package repository

import (
	"abchain_scan/chain"
	"abchain_scan/repository/orm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func preparePairTest() *PairRepository {
	dsn := "host=localhost user=postgres password=12345678 dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	return NewPairRepository(db)
}

func cleanupPairTest(pairRepository *PairRepository, addresses ...string) {
	for _, address := range addresses {
		pairRepository.DeleteByAddressAndChainId(address)
	}
}

func TestPairRepository_Create(t *testing.T) {
	pairRepository := preparePairTest()
	pair := &orm.Pair{
		Name:     "n1",
		Address:  "0x01",
		Token0:   "0x0a",
		Token1:   "0x0b",
		ChainId:  chain.Id,
		Reserve0: decimal.NewFromInt(1),
		Reserve1: decimal.NewFromInt(2),
	}
	createErr := pairRepository.Create(pair)
	require.NoError(t, createErr)

	pairQueried, getErr := pairRepository.GetByAddressAndChainId(pair.Address)
	require.Nil(t, getErr)
	require.True(t, pair.Equal(pairQueried))
	cleanupPairTest(pairRepository, pair.Address)
}

func TestPairRepository_CreateDup(t *testing.T) {
	pairRepository := preparePairTest()
	pairs := []*orm.Pair{
		{
			Name:     "p1",
			Address:  "0xa1",
			Token0:   "0xa1",
			Token1:   "0xa1",
			ChainId:  chain.Id,
			Reserve0: decimal.NewFromInt(1),
			Reserve1: decimal.NewFromInt(2),
		},
		{
			Name:     "p2",
			Address:  "0xa2",
			Token0:   "0xa2",
			Token1:   "0xa2",
			ChainId:  chain.Id,
			Reserve0: decimal.NewFromInt(1),
			Reserve1: decimal.NewFromInt(2),
		},
		{
			Name:     "p3",
			Address:  "0xa3",
			Token0:   "0xa3",
			Token1:   "0xa3",
			ChainId:  chain.Id,
			Reserve0: decimal.NewFromInt(1),
			Reserve1: decimal.NewFromInt(2),
		},
	}

	defer cleanupPairTest(pairRepository, pairs[0].Address, pairs[1].Address, pairs[2].Address)

	createErr := pairRepository.Create(pairs[0])
	require.NoError(t, createErr)

	createBatchErr := pairRepository.CreateBatch(pairs, "address", "chain_id")
	require.NoError(t, createBatchErr)

	for i, pair := range pairs {
		pairQueried, err := pairRepository.GetByAddressAndChainId(pair.Address)
		require.Nil(t, err)
		require.True(t, pairs[i].Equal(pairQueried))
	}
}
