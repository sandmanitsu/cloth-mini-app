package category

import (
	"cloth-mini-app/internal/domain"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/postgresql"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
)

type CategoryRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewCategoryRepository(logger *slog.Logger, db *postgresql.Storage) *CategoryRepository {
	return &CategoryRepository{
		db:     db.DB,
		logger: logger,
	}
}

func (c *CategoryRepository) Categories() ([]domain.Category, error) {
	const op = "repository.Category.Category"

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, _, err := psql.Select("id", "type", "name").From("category").ToSql()
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s : building sql query", op), sl.Err(err))
	}

	rows, err := c.db.Query(sql)
	if err != nil {
		c.logger.Error(fmt.Sprintf("%s: %s", op, sql), sl.Err(err))

		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.CategoryId, &category.Type, &category.Name); err != nil {
			c.logger.Error(op, sl.Err(err))

			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
