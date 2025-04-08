package postgresql

import (
	"context"
	"database/sql"
	"fmt"
)

type txKey string

var (
	ErrGetTransaction = fmt.Errorf("error: getting transaction from context")

	// tx key in context
	ctxTxKey = txKey("tx")
)

func WrapTx(ctx context.Context, db *sql.DB, process func(context.Context) error) error {
	txStarted := txExist(ctx)
	fmt.Println("transact exist ", txStarted)
	if txStarted {
		err := process(ctx)
		if err != nil {
			return err
		}
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx = context.WithValue(ctx, ctxTxKey, tx)

	err = process(ctx)
	fmt.Println("wrapper ", err)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func txExist(ctx context.Context) bool {
	_, ok := ctx.Value(ctxTxKey).(*sql.Tx)

	return ok
}

func TxFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(ctxTxKey).(*sql.Tx)
	if !ok {
		return nil, false
	}

	return tx, true
}
