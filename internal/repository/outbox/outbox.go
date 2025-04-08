package outbox

import (
	domain "cloth-mini-app/internal/domain/event"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
)

type OutboxRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewOutboxRepository(logger *slog.Logger, db *postgresql.Storage) *OutboxRepository {
	return &OutboxRepository{
		db:     db.DB,
		logger: logger,
	}
}

func (o *OutboxRepository) SaveEvent(ctx context.Context, event domain.Event, action func() error) error {
	const op = "storage.postgresql.SaveEvent"

	err := postgresql.WrapTx(ctx, o.db, func(ctx context.Context) error {
		if err := action(); err != nil {
			return err
		}

		if err := o.saveEvent(ctx, event); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		o.logger.Error(fmt.Sprintf("%s : failed create item", op), sl.Err(err))

		return err
	}

	return nil
}

func (o *OutboxRepository) saveEvent(ctx context.Context, event domain.Event) error {
	const op = "storage.postgresql.saveEvent"

	tx, ok := postgresql.TxFromCtx(ctx)
	if !ok {
		o.logger.Error(fmt.Sprintf("%s : failed get transaction from context", op))

		return postgresql.ErrGetTransaction
	}

	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("outbox").
		Columns("event_type", "payload").
		Values(event.EventType, event.Payload).
		ToSql()
	if err != nil {
		o.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return err
	}

	_, err = tx.Exec(sql, args...)
	if err != nil {
		o.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return err
	}

	return nil
}
