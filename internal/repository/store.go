package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Store owns the *gorm.DB and provides a typed transaction helper.
type Store struct{ 
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store { 
	return &Store{
		db: db,
	} 
}
func (s *Store) DB() *gorm.DB    { 
	return s.db 
}

// RunInTx executes fn inside a GORM transaction.
// GORM automatically rolls back on panic or returned error, commits on nil.
func (s *Store) RunInTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := fn(tx); err != nil {
			return fmt.Errorf("tx: %w", err)
		}
		return nil
	})
}