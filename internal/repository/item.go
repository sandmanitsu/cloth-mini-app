package repository

import (
	"cloth-mini-app/internal/domain"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
)

const (
	limitMax = 100 // max records per query if limit doesn't specified
)

type ItemRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewItemRepository(logger *slog.Logger, db *postgresql.Storage) *ItemRepository {
	return &ItemRepository{
		db:     db.DB,
		logger: logger,
	}
}

func (i *ItemRepository) Items(filter map[string]any, limit, offset uint64) ([]domain.ItemAPI, error) {
	const op = "repository.item.Items"

	if limit == 0 {
		limit = limitMax
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	q := psql.Select("i.id", "i.brand", "i.name", "i.description", "i.sex", "i.price", "i.discount", "i.outer_link", "c.type", "c.name AS category_name").
		From("items i").
		LeftJoin("category c on i.category_id = c.id").
		Where(filter).
		Limit(limit).
		Offset(offset)
	sql, args, _ := q.ToSql()

	rows, err := i.db.Query(sql, args...)
	if err != nil {
		i.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return nil, err
	}
	defer rows.Close()

	var items []domain.ItemAPI
	for rows.Next() {
		var item domain.ItemAPI
		if err := rows.Scan(&item.ID, &item.Brand, &item.Name, &item.Description, &item.Sex, &item.Price, &item.Discount, &item.OuterLink, &item.CategoryType, &item.CategoryName); err != nil {
			i.logger.Error(op, sl.Err(err))

			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
