package repository

import (
	domain "cloth-mini-app/internal/domain/lock"
	"cloth-mini-app/internal/storage/postgresql"
	"context"
	"database/sql"
)

type LockRepository struct {
	db *sql.DB
}

func NewLockRepository(db *postgresql.Storage) *LockRepository {
	return &LockRepository{
		db: db.DB,
	}
}

func (l *LockRepository) AdvisoryLock(ctx context.Context, id domain.AdvisoryLockId) error {
	_, err := l.db.Exec("SELECT pg_advisory_lock($1)", id)
	if err != nil {
		return err
	}

	return nil
}

func (l *LockRepository) AdvisoryUnlock(ctx context.Context, id domain.AdvisoryLockId) error {
	_, err := l.db.Exec("SELECT pg_advisory_unlock($1)", id)
	if err != nil {
		return err
	}

	return nil
}
