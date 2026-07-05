package pgsql

import (
	"context"

	"github.com/arthurshafikov/scout-takehome/backend/internal/core/types"
	"gorm.io/gorm"
)

type BaseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) *BaseRepo {
	return &BaseRepo{
		db: db,
	}
}

func (r *BaseRepo) WrapInTransaction(ctx *types.Context, txFunc types.TXFunction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := txFunc(context.WithValue(ctx.GetContext(), types.TXContextKey("tx"), tx)); err != nil {
			return err
		}

		return nil
	})
}

func (r *BaseRepo) getDBInstance(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(types.TXContextKey("tx")).(*gorm.DB); ok {
		return db
	}

	return r.db.WithContext(ctx)
}
