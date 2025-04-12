package storage

import (
	"context"
	"database/sql"
)

type txKey string

var (
	// tx key in context
	ctxTxKey = txKey("tx")
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) getTxFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(ctxTxKey).(*sql.Tx)
	if !ok {
		return nil, false
	}

	return tx, true
}

func (s *Storage) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	tx, ok := s.getTxFromCtx(ctx)
	if !ok {
		return s.db.QueryRow(query, args...)
	}

	return tx.QueryRow(query, args...)
}

func (s *Storage) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	tx, ok := s.getTxFromCtx(ctx)
	if !ok {
		return s.db.Prepare(query)
	}

	return tx.Prepare(query)
}

func (s *Storage) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	tx, ok := s.getTxFromCtx(ctx)
	if !ok {
		return s.db.Query(query, args...)
	}

	return tx.Query(query, args...)
}

func (s *Storage) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx, ok := s.getTxFromCtx(ctx)
	if !ok {
		return s.db.Exec(query, args...)
	}

	return tx.Exec(query, args...)
}
