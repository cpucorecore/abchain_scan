package repository

import (
	"abchain_scan/chain"
	"abchain_scan/repository/orm"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func prepareTokenTest() *TokenRepository {
	dsn := "host=localhost user=postgres password=12345678 dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	return NewTokenRepository(db)
}

func cleanupTokenTest(tokenRepository *TokenRepository, addresses ...string) {
	for _, address := range addresses {
		tokenRepository.DeleteByAddressAndChainId(address)
	}
}

func TestTokenRepository_Create(t *testing.T) {
	tokenRepository := prepareTokenTest()
	address := "0xa1ED0bD9A4776830c5b7bA004F26427b71152CA5"
	token := &orm.Token{
		Address:     address,
		Name:        "DUEL Token",
		Symbol:      "DUEL",
		Decimal:     18,
		TotalSupply: "2659527283.779538",
		ChainId:     chain.Id,
	}

	tokenRepository.Create(token)
	tokenQueried, getErr := tokenRepository.GetByAddressAndChainId(address)
	require.Nil(t, getErr)
	require.True(t, token.Equal(tokenQueried))
	cleanupTokenTest(tokenRepository, address)
}

func TestTokenRepository_CreateBatch_NoConflict(t *testing.T) {
	tokenRepository := prepareTokenTest()
	tokens := []*orm.Token{
		{
			Address:     "0x01",
			Name:        "n1",
			Symbol:      "s1",
			Decimal:     18,
			TotalSupply: "1",
			ChainId:     chain.Id,
		},
		{
			Address:     "0x02",
			Name:        "n2",
			Symbol:      "s2",
			Decimal:     18,
			TotalSupply: "2",
			ChainId:     chain.Id,
		},
		{
			Address:     "0x03",
			Name:        "n3",
			Symbol:      "s3",
			Decimal:     18,
			TotalSupply: "3",
			ChainId:     chain.Id,
		},
	}

	tokenRepository.CreateBatch(tokens, "address", "chain_id")
	for i, token := range tokens {
		tokenQueried, err := tokenRepository.GetByAddressAndChainId(token.Address)
		require.Nil(t, err)
		require.True(t, tokenQueried.Equal(tokens[i]))
	}
	cleanupTokenTest(tokenRepository, tokens[0].Address, tokens[1].Address, tokens[2].Address)
}

func TestTokenRepository_CreateBatch_ConflictWithDb(t *testing.T) {
	tokenRepository := prepareTokenTest()

	tokens := []*orm.Token{
		{
			Address:     "0x01",
			Name:        "n1",
			Symbol:      "s1",
			Decimal:     18,
			TotalSupply: "1",
			ChainId:     chain.Id,
		},
		{
			Address:     "0x02",
			Name:        "n2",
			Symbol:      "s2",
			Decimal:     18,
			TotalSupply: "2",
			ChainId:     chain.Id,
		},
		{
			Address:     "0x03",
			Name:        "n3",
			Symbol:      "s3",
			Decimal:     18,
			TotalSupply: "3",
			ChainId:     chain.Id,
		},
	}

	tokenRepository.Create(tokens[0])
	tokenRepository.CreateBatch(tokens, "address", "chain_id")

	for i, token := range tokens {
		tokenQueried, err := tokenRepository.GetByAddressAndChainId(token.Address)
		require.Nil(t, err)
		require.True(t, tokenQueried.Equal(tokens[i]))
	}

	cleanupTokenTest(tokenRepository, tokens[0].Address, tokens[1].Address, tokens[2].Address)
}

func TestTokenRepository_UpdateMainPair(t *testing.T) {
	tokenRepository := prepareTokenTest()
	token := &orm.Token{
		Address:     "0x01",
		Name:        "n1",
		Symbol:      "s1",
		Decimal:     18,
		TotalSupply: "1",
		ChainId:     chain.Id,
	}

	tokenRepository.Create(token)
	tokenRepository.UpdateMainPair(token.Address, "0x06")

	tokenQueried, err := tokenRepository.GetByAddressAndChainId(token.Address)
	require.Nil(t, err)
	require.True(t, token.Equal(tokenQueried))
	require.Equal(t, "0x06", tokenQueried.MainPair)
	cleanupTokenTest(tokenRepository, token.Address)
}
