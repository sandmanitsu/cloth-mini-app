package postgresql

import (
	"cloth-mini-app/internal/config"
	"cloth-mini-app/internal/retry"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

const (
	duplicateKeyCode = "23505"

	maxRetry       = 5
	baseDelayRetry = time.Second * 5
	maxDelayRetry  = time.Minute
)

type Storage struct {
	DB *sql.DB
}

// Create postgresql db instanse
func NewPostgreSQL(cfg config.DB) (*Storage, error) {
	const op = "storage.postgresql.New"

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBname,
	)

	var db *sql.DB
	var err error

	retryConfig := retry.RetryConfig{
		MaxRetry:  maxRetry,
		BaseDelay: baseDelayRetry,
		MaxDelay:  maxDelayRetry,
		UsedFrom:  "postgresql",
	}
	err = retry.Retry(context.Background(), retryConfig, func() error {
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			return err
		}

		if err = db.Ping(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{DB: db}, nil
}

func IsDuplicateKeyError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == duplicateKeyCode
	}

	return false
}
