package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *BaseRepository[T]) CreateBatch(entities []*T, conflictColumns ...string) error {
	maxBatchSize := 200

	return r.db.Transaction(func(tx *gorm.DB) error {
		for start := 0; start < len(entities); start += maxBatchSize {
			end := start + maxBatchSize
			if end > len(entities) {
				end = len(entities)
			}
			batch := entities[start:end]

			if len(conflictColumns) > 0 {
				columns := make([]clause.Column, len(conflictColumns))
				for i, col := range conflictColumns {
					columns[i] = clause.Column{Name: col}
				}

				tx = tx.Clauses(clause.OnConflict{
					Columns:   columns,
					DoNothing: true,
				})
			}

			if err := tx.Create(batch).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
