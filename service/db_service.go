package service

import (
	"abchain_scan/repository"
	"abchain_scan/repository/orm"
)

type DBService interface {
	AddTokens(tokens []*orm.Token) error
	AddPairs(pairs []*orm.Pair) error
	AddTxs(txs []*orm.Tx) error
}

type dbService struct {
	tokenRepository *repository.TokenRepository
	pairRepository  *repository.PairRepository
	txRepository    *repository.TxRepository
	enableTokenPair bool
	enableTx        bool
}

func (s *dbService) AddTokens(tokens []*orm.Token) error {
	if !s.enableTokenPair {
		return nil
	}

	return s.tokenRepository.CreateBatch(tokens, "address", "chain_id")
}

func (s *dbService) AddPairs(pairs []*orm.Pair) error {
	if !s.enableTokenPair {
		return nil
	}

	return s.pairRepository.CreateBatch(pairs, "address", "chain_id")
}

func (s *dbService) AddTxs(txs []*orm.Tx) error {
	if !s.enableTx {
		return nil
	}

	return s.txRepository.CreateBatch(txs, "token0_address", "block", "block_index", "tx_index")
}

func NewDBService(
	tokenRepository *repository.TokenRepository,
	pairRepository *repository.PairRepository,
	txRepository *repository.TxRepository,
) DBService {
	return &dbService{
		tokenRepository: tokenRepository,
		pairRepository:  pairRepository,
		txRepository:    txRepository,
		enableTokenPair: tokenRepository != nil && pairRepository != nil,
		enableTx:        txRepository != nil,
	}
}
