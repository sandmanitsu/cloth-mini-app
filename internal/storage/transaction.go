package storage

import (
	"context"
	"database/sql"
)

type TransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (t *TransactionManager) RunInTx(ctx context.Context, action func(ctx context.Context) error) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx = context.WithValue(ctx, ctxTxKey, tx)

	err = action(ctx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
