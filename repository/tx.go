package repository

import (
	"abchain_scan/repository/orm"
	"gorm.io/gorm"
)

type TxRepository struct {
	*BaseRepository[orm.Tx]
}

func NewTxRepository(db *gorm.DB) *TxRepository {
	baseRepo := NewBaseRepository[orm.Tx](db)
	return &TxRepository{BaseRepository: baseRepo}
}

func (r *TxRepository) GetByUniqIndex(token0Address string, block uint64, blockIndex, txIndex uint) (*orm.Tx, error) {
	tx := &orm.Tx{}
	err := r.db.Where("token0_address = ? AND block = ? AND block_index = ? AND tx_index = ?",
		token0Address,
		block,
		blockIndex,
		txIndex).First(tx).Error
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *TxRepository) DeleteById(id string) error {
	tx := &orm.Tx{}
	err := r.db.Where("id = ?", id).Delete(tx).Error
	if err != nil {
		return err
	}
	return nil
}
