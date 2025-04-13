package outbox

import (
	domain "cloth-mini-app/internal/domain/event"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
)

const (
	statusDone = "done"
	statusNew  = "new"
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

func (o *OutboxRepository) GetEvents(ctx context.Context) ([]domain.Event, error) {
	var events []domain.Event
	err := postgresql.WrapTx(ctx, o.db, func(ctx context.Context) error {
		eventsBD, err := o.getEvents(ctx)
		if err != nil {
			return err
		}

		if len(eventsBD) == 0 {
			return nil
		}

		events = eventsBD

		if err = o.reserveEvents(ctx, events); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (o *OutboxRepository) getEvents(ctx context.Context) ([]domain.Event, error) {
	const op = "repository.outbox.getEvents"

	tx, ok := postgresql.TxFromCtx(ctx)
	if !ok {
		o.logger.Error(fmt.Sprintf("%s : failed get transaction from context", op))

		return nil, fmt.Errorf("%s : failed get transaction from context", op)
	}

	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "event_type", "payload", "status", "created_at", "reserved_to").
		From("outbox").
		Where("status = ?", statusNew).
		Where("reserved_to IS NULL OR reserved_to < ?", time.Now()).
		ToSql()
	if err != nil {
		o.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))

		return nil, err
	}

	rows, err := tx.Query(sql, args...)
	if err != nil {
		o.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		if err := rows.Scan(
			&event.Id,
			&event.EventType,
			&event.Payload,
			&event.Status,
			&event.CreatedAt,
			&event.ReservedTo,
		); err != nil {
			o.logger.Error(op, sl.Err(err))

			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (o *OutboxRepository) reserveEvents(ctx context.Context, events []domain.Event) error {
	const op = "repository.outbox.ReserveEvents"

	tx, ok := postgresql.TxFromCtx(ctx)
	if !ok {
		o.logger.Error(fmt.Sprintf("%s : failed get transaction from context", op))

		return fmt.Errorf("%s : failed get transaction from context", op)
	}

	reserveTime := time.Now().Add(time.Minute * 5)
	sql, _, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update("outbox").
		Set("reserved_to", "").
		Where("id = ?").
		ToSql()
	if err != nil {
		o.logger.Error(op, sl.Err(err))
		return err
	}

	fmt.Println(sql)
	stmt, err := tx.Prepare(sql)
	if err != nil {
		o.logger.Error(fmt.Sprintf("%s : failed get prepare stmt", op))

		return err
	}

	for _, event := range events {
		if _, err := stmt.Exec(reserveTime, event.Id); err != nil {
			o.logger.Error(fmt.Sprintf("%s : failed reserve event", op))
			return err
		}
	}

	return nil
}

func (o *OutboxRepository) CreateEvent(ctx context.Context, event domain.Event) error {
	const op = "repository.outbox.CreateEvent"

	tx, ok := postgresql.TxFromCtx(ctx)
	if !ok {
		o.logger.Error(fmt.Sprintf("%s : failed get transaction from context", op))

		return fmt.Errorf("%s : failed get transaction from context", op)
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

func (o *OutboxRepository) ChangeStatus(ctx context.Context, eventsId []int) error {
	const op = "repository.outbox.ChangeStatus"

	sql, _, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update("outbox").
		Set("status", "").
		Where("id = ?").
		ToSql()
	if err != nil {
		o.logger.Error(op, sl.Err(err))
		return err
	}

	fmt.Println(sql)
	stmt, err := o.db.Prepare(sql)
	if err != nil {
		o.logger.Error(fmt.Sprintf("%s : failed get prepare stmt", op))

		return err
	}

	for _, id := range eventsId {
		if _, err := stmt.Exec(statusDone, id); err != nil {
			o.logger.Error(fmt.Sprintf("%s : failed reserve event", op))
			return err
		}
	}

	return nil
}
